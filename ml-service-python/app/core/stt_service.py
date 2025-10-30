import io
import numpy as np
import torch
import soundfile as sf
from transformers import WhisperProcessor, WhisperForConditionalGeneration
from typing import List

class STTService:
    def __init__(self, model_name: str = "openai/whisper-small"):
        self.processor = WhisperProcessor.from_pretrained(model_name)
        self.model = WhisperForConditionalGeneration.from_pretrained(model_name)
        self.model.eval()
        self.device = "cuda" if torch.cuda.is_available() else "cpu"
        self.model.to(self.device)
        print(f"Whisper loaded on {self.device}")

        # VAD: simple energy-based (700ms silence = end)
        self.silence_threshold = 0.01
        self.silence_duration_ms = 700
        self.sample_rate = 16000
        self.chunk_duration_ms = 100  # assume 100ms per chunk

    def add_audio_chunk(self, audio_bytes: bytes, current_buffer: List[np.ndarray]) -> List[np.ndarray]:
        """Convert bytes → PCM → append to buffer"""
        try:
            audio_np, sr = sf.read(io.BytesIO(audio_bytes), dtype='float32')
            if sr != self.sample_rate:
                raise ValueError(f"Expected {self.sample_rate}Hz, got {sr}Hz")
            current_buffer.append(audio_np)
            return current_buffer
        except Exception as e:
            print(f"Audio decode error: {e}")
            return current_buffer

    def detect_silence(self, buffer: List[np.ndarray]) -> bool:
        """Check last N ms for silence"""
        if not buffer:
            return False
        recent = np.concatenate(buffer[-int(self.silence_duration_ms / self.chunk_duration_ms):])
        energy = np.mean(np.abs(recent))
        return energy < self.silence_threshold

    def transcribe_buffer(self, buffer: List[np.ndarray]) -> str:
        """Run Whisper on full buffer"""
        if not buffer:
            return ""
        audio = np.concatenate(buffer)
        input_features = self.processor(
            audio, sampling_rate=self.sample_rate, return_tensors="pt"
        ).input_features.to(self.device)

        with torch.no_grad():
            predicted_ids = self.model.generate(input_features, max_new_tokens=448)
        transcription = self.processor.batch_decode(predicted_ids, skip_special_tokens=True)[0]
        return transcription.strip()

# Global instance
stt_service = STTService()