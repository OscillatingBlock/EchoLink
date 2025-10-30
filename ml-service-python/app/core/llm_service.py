# app/core/llm_service.py
from transformers import pipeline

# Load model once during startup
# Using small instruct model to reduce VPS resource usage
llm = pipeline("text-generation", model="mistralai/Mistral-7B-Instruct-v0.2")

class LLMService:

    @staticmethod
    def generate_reply(history: list) -> str:
        """
        history: [{"role": "user"/"assistant", "text": "..."}]

        Converts chat history → conversation prompt
        → calls model → returns assistant response.
        """

        prompt = ""
        for msg in history:
            prompt += f"{msg['role']}: {msg['text']}\n"
        prompt += "assistant:"

        result = llm(prompt, max_new_tokens=100)[0]["generated_text"]
        reply = result.split("assistant:")[-1].strip()
        return reply


llm_service = LLMService()
