package main

import (
	"context"
	"eino/config"
	"fmt"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	embedding "github.com/cloudwego/eino-ext/components/embedding/openai"
	redisInd "github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino-ext/components/model/ark"
	redisRet "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

// RAGEngine 定义了一个基于检索增强生成（RAG）的引擎结构体
type RAGEngine struct {
	indexName string // 索引名称
	prefix    string // 前缀
	config    *config.ParamsConfig
	dimension int // 向量维度

	redis    *redis.Client
	embedder *embedding.Embedder // 嵌入器，用于将文本转换为向量表示
	Err      error

	Loader    *file.FileLoader     // 用于加载文档的加载器
	Splitter  document.Transformer // 用于拆分文档的拆分器
	Retriever *redisRet.Retriever  // 用于检索相关文档的检索器
	Indexer   *redisInd.Indexer    // 用于索引文档的索引器
	ChatModel *ark.ChatModel       // 用于生成回答的聊天模型
}

func InitRAGEngine(ctx context.Context, index string, prefix string) (*RAGEngine, error) {
	r, err := initRAGEngine(ctx, index, prefix)
	if err != nil {
		return nil, err
	}

	r.newLoader(ctx)    // 初始化加载器, 用于加载文档
	r.newSplitter(ctx)  // 初始化拆分器, 用于拆分文档
	r.newIndexer(ctx)   // 创建 Redis 文档索引
	r.newRetriever(ctx) // 初始化检索器, 用于检索文档
	r.newChatModel(ctx) // 初始化聊天模型, 用于生成回答

	return r, nil
}

func initRAGEngine(ctx context.Context, index string, prefix string) (*RAGEngine, error) {

	c := config.Map()

	// 创建 embedder 用于将文档转成向量，Indexer与Retriever均依赖于该组件。
	embedder, err := embedding.NewEmbedder(ctx, &embedding.EmbeddingConfig{
		// APIKey: c.ApiKey,    // API Key
		Model:   c.Embedding, // 使用的模型名称
		BaseURL: c.EmbeddingBaseURL,
	})

	if err != nil {
		return nil, err
	}

	return &RAGEngine{
		indexName: index,
		prefix:    prefix,
		config:    c,
		dimension: 4096,

		redis: redis.NewClient(&redis.Options{
			Addr:          fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
			Protocol:      2,
			UnstableResp3: true,
		}),
		embedder: embedder,

		Loader:    nil,
		Splitter:  nil,
		Retriever: nil,
		Indexer:   nil,
		ChatModel: nil,
	}, nil
}

// systemPrompt 定义了系统提示模板，指导模型如何生成回答
var systemPrompt = `
# Role: Student Learning Assistant

# Language: Chinese

- When providing assistance:
  • Be clear and concise
  • Include practical examples when relevant
  • Reference documentation when helpful
  • Suggest improvements or next steps if applicable

here's documents searched for you:
==== doc start ====
	  {documents}
==== doc end ====
`

// Generate 方法根据用户的查询生成回答
func (r *RAGEngine) Generate(ctx context.Context, query string) (*schema.StreamReader[*schema.Message], error) {

	// 检索相关文档
	docs, err := r.Retriever.Retrieve(ctx, query)
	if err != nil {
		return nil, err
	}

	fmt.Println("-------------------------------------------")
	fmt.Println(docs)
	fmt.Println("-------------------------------------------")

	// 格式化提示模板
	tpl := prompt.FromMessages(schema.FString, []schema.MessagesTemplate{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage("question: {content}"),
	}...)

	messages, err := tpl.Format(ctx, map[string]any{
		"documents": docs,  // 将检索到的文档作为上下文
		"content":   query, // 将用户查询作为问题
	})
	if err != nil {
		return nil, err
	}

	// 生成回答
	return r.ChatModel.Stream(ctx, messages)
}
