package elastic

import (
	"context"
	"fmt"

	"github.com/liyiysng/scatter/config"
	"github.com/olivere/elastic/v7"
)

// SetupEnv 初始化环境
func SetupEnv(client *elastic.Client, cfg *config.Config) error {

	templateName := cfg.GetString("scatter.es.template_name")

	// 判断模板是否存在
	te := client.IndexTemplateExists(templateName)
	teExist, err := te.Do(context.Background())
	if err != nil {
		return err
	}

	// 环境已经设置
	if teExist {
		return nil
	}

	err = createLifecycle(client, cfg)
	if err != nil {
		return err
	}

	err = createTemplate(client, cfg)
	if err != nil {
		return err
	}

	// 创建默认的
	initIndexName := cfg.GetString("scatter.es.init_index_name")
	initIndex := client.CreateIndex(initIndexName)

	strBody := fmt.Sprintf(`{"aliases":{ "%s":{ "is_write_index": true } } }`, cfg.GetString("scatter.es.template.settings.index.lifecycle.rollover_alias"))

	initIndex.BodyJson(strBody)

	_, err = initIndex.Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// 创建生命周期管理
func createLifecycle(client *elastic.Client, cfg *config.Config) error {
	lp := client.XPackIlmPutLifecycle()

	lp.Policy(cfg.GetString("scatter.es.lifecycle_name"))

	lp.BodyJson(cfg.Get("scatter.es.lifecycle"))

	lpRes, err := lp.Do(context.Background())
	if err != nil {
		return err
	}
	if lpRes.Acknowledged {
		myLog.Infof("create lifecycle %s success", cfg.GetString("scatter.es.lifecycle_name"))
	}
	return nil
}

// 创建模板
func createTemplate(client *elastic.Client, cfg *config.Config) error {

	name := cfg.GetString("scatter.es.template_name")

	tp := client.IndexPutTemplate(name)

	//tp := client.IndexPutIndexTemplate(name)
	tp.BodyJson(cfg.Get("scatter.es.template"))

	tpRes, err := tp.Do(context.Background())
	if err != nil {
		return err
	}
	myLog.Infof("create template %s res %v", name, tpRes)
	return nil
}
