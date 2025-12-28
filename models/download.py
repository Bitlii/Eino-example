#!/usr/bin/env python3

from transformers import AutoTokenizer, AutoModelForCausalLM

model_name = "Qwen/Qwen3-Embedding-0.6B"
save_dir = "./Qwen3-Embedding-0.6B"

# 直接下载并加载模型（会自动缓存）
tokenizer = AutoTokenizer.from_pretrained(model_name, cache_dir="./cache")
model = AutoModelForCausalLM.from_pretrained(model_name, cache_dir="./cache")

tokenizer.save_pretrained(save_dir)
model.save_pretrained(save_dir)

print(f"模型已成功下载并保存到 {save_dir}")