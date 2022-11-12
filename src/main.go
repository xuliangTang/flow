package main

import (
	"flow/src/athena"
	"flow/src/classes"
	"flow/src/conf"
)

func Div(v1 int, v2 int) int {
	if v2 == 0 {
		return 0
	}
	return v1 / v2
}

func main() {
	/*cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"39.103.214.57:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer cli.Close()

	kv := clientv3.NewKV(cli)
	ctx := context.Background()
	_, err = kv.Put(ctx, "/service/test3", "test3-service")
	if err != nil {
		log.Fatalln(err)
	}*/

	/*r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong v2",
		})
	})
	r.Run(":8080")*/

	athena.Ignite().
		Load(conf.NewConfigModule()).
		Mount("v1", nil, classes.NewIndexClass()).
		Launch()
}
