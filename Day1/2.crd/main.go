package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func generateRandomString(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)[:n]
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "custom.example.com",
		Version:  "v1",
		Resource: "secretgenerators",
	}

	crds, err := dynamicClient.Resource(gvr).Namespace("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, cr := range crds.Items {
		secretName, found, err := unstructured.NestedString(cr.Object, "spec", "secretName")
		if err != nil || !found {
			fmt.Println("Secret name not found in CRD")
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
}