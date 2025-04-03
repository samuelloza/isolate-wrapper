package services_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/samuelloza/isolate-wrapper/src/application/services"
)

// TestIsSandboxIDFree verifica que la función IsSandboxIDFree retorne true cuando el sandbox está libre
// y false cuando el directorio existe.
func TestIsSandboxIDFree(t *testing.T) {
	manager := services.NewSandboxManagerService(10)
	testBoxID := 42
	testPath := fmt.Sprintf("/var/local/lib/isolate/%d", testBoxID)

	// Asegurarse de que el sandbox esté libre: eliminar el directorio si existe.
	_ = os.RemoveAll(testPath)

	// Si no existe el directorio, se espera que IsSandboxIDFree retorne true.
	if !manager.IsSandboxIDFree(testBoxID) {
		t.Errorf("Se esperaba que el sandbox ID %d estuviera libre", testBoxID)
	}

	// Crear el directorio para simular que el sandbox está ocupado.
	if err := os.MkdirAll(testPath, 0755); err != nil {
		t.Fatalf("Error al crear el directorio de prueba: %v", err)
	}
	// Asegurarse de limpiar luego de la prueba.
	defer os.RemoveAll(testPath)

	// Ahora se espera que IsSandboxIDFree retorne false.
	if manager.IsSandboxIDFree(testBoxID) {
		t.Errorf("Se esperaba que el sandbox ID %d estuviera ocupado", testBoxID)
	}
}

// TestGetAvailableSandboxID verifica que, al iniciar desde un sandbox ocupado,
// se retorne otro ID disponible.
func TestGetAvailableSandboxID(t *testing.T) {
	// Se utiliza un número pequeño de cajas para simplificar la prueba.
	maxBoxes := 5
	manager := services.NewSandboxManagerService(maxBoxes)
	initialBoxID := 1

	// Simular que el sandbox con ID inicial está ocupado creando su directorio.
	occupiedPath := fmt.Sprintf("/var/local/lib/isolate/%d", initialBoxID)
	if err := os.MkdirAll(occupiedPath, 0755); err != nil {
		t.Fatalf("Error al crear el directorio ocupado: %v", err)
	}
	defer os.RemoveAll(occupiedPath)

	// Se llama a GetAvailableSandboxID con el ID ocupado.
	startTime := time.Now()
	boxID, err := manager.GetAvailableSandboxID(initialBoxID)
	elapsed := time.Since(startTime)

	if err != nil {
		t.Fatalf("Error al obtener un sandbox disponible: %v", err)
	}
	if boxID == initialBoxID {
		t.Errorf("Se esperaba un sandbox ID distinto a %d, pero se obtuvo el mismo", initialBoxID)
	}
	// Se puede advertir si la función no tardó lo esperado (mínimo 3 segundos por el sleep).
	if elapsed < 3*time.Second {
		t.Logf("Advertencia: Se esperaba un retardo de al menos 3 segundos, pero se tardó %v", elapsed)
	}
}

// (Opcional) Test para cuando no hay ningún sandbox disponible.
func TestGetAvailableSandboxID_AllOccupied(t *testing.T) {
	maxBoxes := 2
	manager := services.NewSandboxManagerService(maxBoxes)

	// Ocupamos todos los sandbox IDs disponibles (0 y 1).
	for i := 0; i < maxBoxes; i++ {
		path := fmt.Sprintf("/var/local/lib/isolate/%d", i)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("Error al ocupar el sandbox ID %d: %v", i, err)
		}
		defer os.RemoveAll(path)
	}

	_, err := manager.GetAvailableSandboxID(0)
	if err == nil {
		t.Errorf("Se esperaba error al no haber sandbox disponibles, pero no se obtuvo ninguno")
	}
}
