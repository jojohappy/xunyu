package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func genConfig(path string, reOpen bool, follow bool) FileConfig {
	return FileConfig{
		Path:   path,
		ReOpen: reOpen,
		Follow: follow,
	}
}

func TestGenTailConfig(t *testing.T) {
	path, _ := filepath.Abs(filepath.Join("../../..", "testdata", "data.log"))
	filConfig := genConfig(path, false, true)
	f := &file{config: filConfig}

	_, err := f.newTailConfig()

	if nil != err {
		t.Error(err)
		t.Fail()
	}
}

func TestTailFile(t *testing.T) {
	path, _ := filepath.Abs(filepath.Join("../../..", "testdata", "data.log"))
	filConfig := genConfig(path, false, true)
	f := &file{config: filConfig}
	count := 0
	createFile(path, "{\"real_ip\":\"-\",\"remote\":\"218.26.54.2\",\"port\":\"6232\",\"host\":\"example.app.com\"}\n{\"real_ip\":\"-\",\"remote\":\"180.255.250.216\",\"port\":\"52086\",\"host\":\"example.test.com\"}\n{\"real_ip\":\"-\",\"remote\":\"110.17.170.218\",\"port\":\"38768\",\"host\":\"example.bbs.com\"}\n{\"real_ip\":\"-\",\"remote\":\"223.104.3.236\",\"port\":\"41206\",\"host\":\"example.a.com\"}\n{\"real_ip\":\"-\",\"remote\":\"223.104.16.81\",\"port\":\"35357\",\"host\":\"example.b.com\"}\n", t)

	go func() {
		out := f.Start()
		for data := range out {
			t.Logf("test input plugin file, received message: %v", data)
			count++
		}
	}()
	<-time.After(time.Second * 1)
	t.Logf("test input plugin file, total received messages: %d", count)

	count = 0
	appendFile(path, "{\"real_ip\":\"-\",\"remote\":\"218.24.121.4\",\"port\":\"19725\",\"host\":\"example.cc.com\"}\n", t)
	<-time.After(time.Second)
	t.Logf("test input plugin file, total received messages: %d", count)
	removeFile(path, t)
}

func appendFile(path string, contents string, t *testing.T) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(contents)
	if err != nil {
		t.Fatal(err)
	}
}

func createFile(path string, contents string, t *testing.T) {
	err := ioutil.WriteFile(path, []byte(contents), 0600)
	if err != nil {
		t.Fatal(err)
	}
}

func removeFile(path string, t *testing.T) {
	err := os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}
}
