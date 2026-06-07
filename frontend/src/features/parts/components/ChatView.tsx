import React, { useState, useEffect } from 'react';
import { Send, Bot, User, ArrowLeft } from 'lucide-react';
import { chatApi } from "../services/partsApi";

interface Message {
  id: string;
  sender: 'bot' | 'user';
  text: string;
  time: string;
}

interface ChatViewProps {
  onBackToGarage: () => void;
}

export default function ChatView({ onBackToGarage }: ChatViewProps) {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: 'msg-1',
      sender: 'bot',
      text: '¡Hola! Soy el asistente del taller. ¿En qué puedo ayudarte?',
      time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    }
  ]);
  const [inputValue, setInputValue] = useState('');
  const [sessionId, setSessionId] = useState<number | null>(null);

  // 1. Al montar el componente, inicializamos la sesión con el backend
  useEffect(() => {
    const initSession = async () => {
      try {
        const response = await chatApi.createSession({
          vehicle_brand: "Audi", 
          vehicle_model: "A4", 
          vehicle_year: "2024"
        });
        setSessionId(response.id);
      } catch (error) {
        console.error("Error al crear sesión:", error);
      }
    };
    initSession();
  }, []);

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim() || !sessionId) return;

    const userText = inputValue;
    const userTime = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    
    // A. Añadimos el mensaje del usuario inmediatamente a la UI
    const newUserMessage: Message = {
      id: `msg-user-${Date.now()}`,
      sender: 'user',
      text: userText,
      time: userTime
    };

    setMessages(prev => [...prev, newUserMessage]);
    setInputValue('');

    // B. Llamamos a TU Backend (Go)
    try {
      const response = await chatApi.sendMessage(sessionId, userText);
      
      // C. Añadimos la respuesta del bot (que viene de n8n -> tu Go)
      const botTime = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      setMessages(prev => [
        ...prev,
        {
          id: `msg-bot-${Date.now()}`,
          sender: 'bot',
          text: response.assistantMessage.content, // Aquí viene la respuesta real
          time: botTime
        }
      ]);
    } catch (error) {
      console.error("Error enviando mensaje:", error);
      setMessages(prev => [...prev, {
        id: `msg-err-${Date.now()}`,
        sender: 'bot',
        text: "Error: No pude conectar con el servidor.",
        time: new Date().toLocaleTimeString()
      }]);
    }
  };

  return (
    <div className="w-full h-[calc(100vh-72px)] bg-slate-950 flex flex-col">
      {/* Barra superior */}
      <div className="px-6 py-3 bg-slate-900/40 backdrop-blur-md border-b border-slate-800 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button onClick={onBackToGarage} className="p-1.5 rounded-lg bg-slate-950 border border-slate-800 text-slate-400 hover:text-cyan-400 hover:border-cyan-500/50 transition-all cursor-pointer">
            <ArrowLeft className="w-4 h-4" />
          </button>
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 bg-cyan-950 border border-cyan-500/30 rounded-full flex items-center justify-center relative">
              <Bot className="text-cyan-400 w-5 h-5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-slate-200">Asistente Taller</h3>
              <p className="text-[10px] text-emerald-400 font-medium">Conectado a Backend</p>
            </div>
          </div>
        </div>
      </div>

      {/* Cuerpo de Mensajes */}
      <div className="flex-1 overflow-y-auto p-6 flex flex-col gap-4">
        {messages.map((msg) => (
          <div key={msg.id} className={`flex gap-3 max-w-[75%] ${msg.sender === 'bot' ? 'self-start' : 'self-end flex-row-reverse'}`}>
            <div className={`w-7 h-7 rounded-full flex items-center justify-center flex-shrink-0 ${msg.sender === 'bot' ? 'bg-slate-900 border border-slate-800 text-cyan-400' : 'bg-cyan-950 border border-cyan-800 text-cyan-300'}`}>
              {msg.sender === 'bot' ? <Bot className="w-3.5 h-3.5" /> : <User className="w-3.5 h-3.5" />}
            </div>
            <div className="flex flex-col gap-1">
              <div className={`rounded-xl p-3 text-xs leading-relaxed border ${msg.sender === 'bot' ? 'bg-slate-900/60 border-slate-800 text-slate-300' : 'bg-gradient-to-br from-cyan-950/60 to-slate-900/60 border-cyan-900/50 text-cyan-100'}`}>
                {msg.text}
              </div>
              <span className="text-[9px] text-slate-500 font-mono">{msg.time}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Formulario */}
      <form onSubmit={handleSendMessage} className="p-4 bg-slate-900/20 border-t border-slate-900 flex gap-3 items-center">
        <input 
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          placeholder="Pregunta sobre compatibilidad..."
          className="flex-1 bg-slate-950 border border-slate-800 rounded-lg px-4 py-2.5 text-xs text-slate-300 focus:outline-none focus:border-cyan-500/50"
        />
        <button type="submit" disabled={!inputValue.trim() || !sessionId} className="p-2.5 bg-cyan-950 text-cyan-400 border border-cyan-800 hover:bg-cyan-400 hover:text-slate-950 rounded-lg disabled:opacity-40">
          <Send className="w-4 h-4" />
        </button>
      </form>
    </div>
  );
}