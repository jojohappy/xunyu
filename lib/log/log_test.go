package log

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/xunyu/config"
)

type node map[string]interface{}

func genConfigNoLevel() (*config.Config, error) {
	cfg, err := config.From(node{
		"file": "./logs/xunyu_nolevel.log",
	})

	return cfg, err
}

func genConfigNoFile() (*config.Config, error) {
	cfg, err := config.From(node{
		"level": 4,
	})

	return cfg, err
}

func TestLoadLogConfig(t *testing.T) {
	cfg, _ := config.From(node{})

	err := InitLog(cfg)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("instance _log is %v\n", _log)

	if _, err := os.Stat(_log.config.File); os.IsNotExist(err) {
		t.Error(err)
	}

	os.Remove(_log.config.File)
	os.Remove(path.Dir(_log.config.File))
}

func TestLogNoLevel(t *testing.T) {
	cfg, _ := genConfigNoLevel()
	testLog(cfg, t)
}

func TestLogNoFile(t *testing.T) {
	cfg, _ := genConfigNoFile()
	testLog(cfg, t)
}

func testLog(cfg *config.Config, t *testing.T) {
	err := InitLog(cfg)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("instance _log is %v\n", _log)

	if _, err := os.Stat(_log.config.File); os.IsNotExist(err) {
		t.Error(err)
	}

	Error("test %s", "test1")

	content, err := ioutil.ReadFile(_log.config.File)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("file content is\n%s\n", string(content))

	Debug("test debug %s", "test debug")

	content, err = ioutil.ReadFile(_log.config.File)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("file content after debug is\n%s\n", string(content))

	os.Remove(_log.config.File)
	os.Remove(path.Dir(_log.config.File))
}
