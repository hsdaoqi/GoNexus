import chromadb
from chromadb.utils import embedding_functions
from openai import OpenAI
from core.config import settings


class RAGEngine:
    def __init__(self):
        print("🚀 [Init] 正在启动 RAG 引擎...", flush=True)

        # 1. 连接 AI
        self.ai_client = OpenAI(
            api_key=settings.API_KEY,
            base_url=settings.BASE_URL
        )

        # 2. 连接数据库
        print("--- 正在初始化数据库 (首次运行需下载模型，请稍候)... ---", flush=True)
        self.chroma_client = chromadb.PersistentClient(path=settings.DB_PATH)

        # 3. 加载模型
        emb_fn = embedding_functions.DefaultEmbeddingFunction()

        self.collection = self.chroma_client.get_or_create_collection(
            name="chat_memory",
            embedding_function=emb_fn
        )
        print("✅ [Init] RAG 引擎启动成功！", flush=True)

    # 存入记忆 (带昵称)
    def add_memory(self, text: str, user_id: str, msg_id: str, nickname: str, session_id: str):
        print(f"📥 [记忆] 会话:{session_id} {nickname}: {text}")
        self.collection.add(
            documents=[text],
            metadatas=[{"user_id": user_id, "user": nickname, "session_id": session_id}],
            ids=[msg_id]
        )

    # 检索记忆
    def search_memory(self, query: str, session_id: str, limit: int = 3):
        print(f"🔍 [检索] 会话:{session_id} 问题:{query}")
        results = self.collection.query(
            query_texts=[query],
            n_results=limit,
            # 核心安全锁：只搜这个 session_id 的数据
            where={"session_id": session_id}
        )

        if not results['documents'] or not results['documents'][0]:
            return []

        docs = results['documents'][0]
        metas = results['metadatas'][0]
        combined = []

        for i in range(len(docs)):
            # 优先取昵称
            name = metas[i].get('user', metas[i].get('user_id', '未知'))
            text = docs[i]
            combined.append(f"{name}: {text}")

        return combined

    # RAG 问答
    def rag_qa(self, question: str, session_id: str, limit: int):
        # 1. 检索 (把 limit 改大一点，比如 10 条)
        # 既然是聊天记录，上下文多一点没坏处
        related_docs = self.search_memory(question, session_id, limit)

        # 🔥🔥🔥【调试代码】打印出 AI 到底看到了什么
        print(f"🧐 [Debug] 用户问: {question}")
        print(f"🧐 [Debug] 检索到的上下文: {related_docs}")
        # 🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥

        if not related_docs:
            return "抱歉，我的记忆库里没有关于这件事的记录。"

        context = "\n".join(related_docs)

        # 2. 生成 (Prompt 微调：让 AI 更聪明一点)
        prompt = f"""
        你是一个聊天记录分析助手。请根据【参考记录】回答【用户问题】。

        注意：
        1. 记录格式为 "姓名: 内容"。
        2. "我" 指的是发言者。例如 "Alice: 我饿了" 意味着 Alice 饿了。
        3. 请根据上下文逻辑推理。

        【参考记录】：
        {context}

        【用户问题】：{question}
        """

        # 3. 提问 AI
        print(f"🤖 [思考] 正在请求 AI...")
        response = self.ai_client.chat.completions.create(
            model=settings.MODEL_NAME,
            messages=[{"role": "user", "content": prompt}]
        )
        return response.choices[0].message.content

    # 总结
    def chat_summary(self, chats: list):
        context = "\n".join(chats)
        prompt = f"请总结以下聊天内容：\n{context}"
        response = self.ai_client.chat.completions.create(
            model=settings.MODEL_NAME,
            messages=[{"role": "user", "content": prompt}]
        )
        return response.choices[0].message.content
