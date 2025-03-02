package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func generateRandomString(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println("Error generating random string:", err)
		return ""
	}
	return hex.EncodeToString(bytes)[:n]
}

func main() {
	fmt.Println("Starting secret-operator CronJob...")

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Error getting in-cluster config:", err)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		return
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating dynamic client:", err)
		return
	}

	secretGenGVR := schema.GroupVersionResource{
		Group:    "custom.example.com",
		Version:  "v1",
		Resource: "secretgenerators",
	}

	crs, err := dynamicClient.Resource(secretGenGVR).Namespace("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error fetching SecretGenerators:", err)
		return
	}

	for _, cr := range crs.Items {
		secretName, found, err := unstructuredNestedString(cr.Object, "spec", "secretName")
		if err != nil || !found {
			fmt.Println("SecretGenerator missing 'spec.secretName', skipping:", cr.GetName())
			continue
		}

		existingSecret, err := clientset.CoreV1().Secrets("default").Get(context.TODO(), secretName, metav1.GetOptions{})
		if err == nil {
			existingSecret.Data["password"] = []byte(generateRandomString(8))
			_, err = clientset.CoreV1().Secrets("default").Update(context.TODO(), existingSecret, metav1.UpdateOptions{})
			if err != nil {
				fmt.Println("Error updating secret:", err)
			} else {
				fmt.Println("Updated secret:", secretName)
			}
		} else {
			newSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretName,
					Namespace: "default",
					Labels:    map[string]string{"managed-by": "secret-operator"},
				},
				Data: map[string][]byte{"password": []byte(generateRandomString(8))},
			}
			_, err := clientset.CoreV1().Secrets("default").Create(context.TODO(), newSecret, metav1.CreateOptions{})
			if err != nil {
				fmt.Println("Error creating secret:", err)
			} else {
				fmt.Println("Created secret:", secretName)
			}
		}
	}

	fmt.Println("Secret-operator CronJob finished.")
}

func unstructuredNestedString(obj map[string]interface{}, fields ...string) (string, bool, error) {
	val, found, err := unstructured.NestedString(obj, fields...)
	return val, found, err
}