package main

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/eino/components/document"
	uuid2 "github.com/google/uuid"
)

const (
	prefix = "OuterCyrex:" // 前缀
	index  = "OuterIndex"  // 索引名称
)

func main() {
	ctx := context.Background()

	// 通过 indexName 和 prefix 初始化 RAGEngine
	r, err := InitRAGEngine(ctx, index, prefix)
	if err != nil {
		panic(err)
	}

	// 加载一个文档
	doc, err := r.Loader.Load(ctx, document.Source{
		URI: "./test_txt/mysql-1.md",
	})
	if err != nil {
		panic(err)
	}

	// 拆分文档，一个被拆分成多个小文档
	docs, err := r.Splitter.Transform(ctx, doc)
	if err != nil {
		panic(err)
	}

	// 为每个文档生成唯一 ID
	for _, d := range docs {
		uuid, _ := uuid2.NewUUID()
		d.ID = uuid.String()
	}

	// 初始化向量索引， Redis 中创建索引
	err = r.InitVectorIndex(ctx)
	if err != nil {
		panic(err)
	}

	// 将文档索引到向量数据库中
	// 目前发现多次运行会多次存入数据，导致重复数据增多，后续需要优化
	_, err = r.Indexer.Store(ctx, docs)
	if err != nil {
		panic(err)
	}

	var query string

	for {
		// 等待输入
		_, _ = fmt.Scan(&query)
		output, err := r.Generate(ctx, query)
		if err != nil {
			panic(err)
		}
		for {
			o, err := output.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			fmt.Println(o.Content)
		}
	}
}
