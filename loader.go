package main

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
)

func (r *RAGEngine) newLoader(ctx context.Context) {
	// 创建文件加载器
	l, err := file.NewFileLoader(ctx, &file.FileLoaderConfig{
		UseNameAsID: true, // 使用文件名作为文档ID
		Parser:      nil,  // 使用默认的文件解析器
	})
	if err != nil {
		r.Err = err
		return
	}
	r.Loader = l
}
