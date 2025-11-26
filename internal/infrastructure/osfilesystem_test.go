package infrastructure

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestOSFileSystem_WalkDir(t *testing.T) {
	fileSys := &OSFileSystem{}
	tempDir := t.TempDir()
	// Create some files
	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("content1"), 0644)
	os.MkdirAll(filepath.Join(tempDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tempDir, "subdir", "file2.txt"), []byte("content2"), 0644)

	var walked []string
	err := fileSys.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		walked = append(walked, path)
		return nil
	})
	if err != nil {
		t.Errorf("WalkDir failed: %v", err)
	}
	if len(walked) < 3 { // tempDir, file1, subdir, file2
		t.Error("WalkDir did not walk all files")
	}
}

func TestOSFileSystem_ReadFile(t *testing.T) {
	fileSys := &OSFileSystem{}
	tempFile := filepath.Join(t.TempDir(), "test.txt")
	content := []byte("test content")
	os.WriteFile(tempFile, content, 0644)

	readContent, err := fileSys.ReadFile(tempFile)
	if err != nil {
		t.Errorf("ReadFile failed: %v", err)
	}
	if string(readContent) != string(content) {
		t.Error("ReadFile returned wrong content")
	}
}

func TestOSFileSystem_Create(t *testing.T) {
	fileSys := &OSFileSystem{}
	tempFile := filepath.Join(t.TempDir(), "test.txt")
	f, err := fileSys.Create(tempFile)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}
	defer f.Close()
	f.Write([]byte("test"))
	// Check file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("File not created")
	}
}

func TestOSFileSystem_MkdirAll(t *testing.T) {
	fileSys := &OSFileSystem{}
	tempDir := filepath.Join(t.TempDir(), "testdir", "subdir")
	err := fileSys.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Errorf("MkdirAll failed: %v", err)
	}
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Directory not created")
	}
}

func TestOSFileSystem_RemoveAll(t *testing.T) {
	fileSys := &OSFileSystem{}
	tempDir := filepath.Join(t.TempDir(), "testdir")
	os.MkdirAll(tempDir, 0755)
	os.WriteFile(filepath.Join(tempDir, "file.txt"), []byte(""), 0644)

	err := fileSys.RemoveAll(tempDir)
	if err != nil {
		t.Errorf("RemoveAll failed: %v", err)
	}
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Directory not removed")
	}
}
