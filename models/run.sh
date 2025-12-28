#!/bin/bash

python -m vllm.entrypoints.openai.api_server \
    --model ./Qwen3-Embedding-0.6B \
    --served-model-name Qwen3-Embedding-0.6B \
    --enforce-eager \
    --max-model-len 32768 \
    --port 8000 \
    --host 127.0.0.1 \
    --convert embed \
    --runner pooling