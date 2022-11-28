package setting

import (
	"testing"

	"github.com/zhangbiao651/fabric-manager/web/models"
	"github.com/zhangbiao651/fabric-manager/web/pkg/setting"
)

func TestSetup(t *testing.T) {
	setting.Setup()
	models.Setup()
}
