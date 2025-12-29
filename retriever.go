package main

import (
	"context"

	redisRet "github.com/cloudwego/eino-ext/components/retriever/redis"
)

func (r *RAGEngine) newRetriever(ctx context.Context) {
	// Redis 检索器
	re, err := redisRet.NewRetriever(ctx, &redisRet.RetrieverConfig{
		Client:            r.redis,          // Redis 客户端
		Index:             r.indexName,      // 索引名称
		VectorField:       "vector_content", // 向量字段
		DistanceThreshold: nil,
		Dialect:           2,
		ReturnFields:      []string{"vector_content", "content"}, // 返回字段
		DocumentConverter: nil,
		TopK:              1,
		Embedding:         r.embedder,
	})

	if err != nil {
		r.Err = err
		return
	}

	r.Retriever = re
}
