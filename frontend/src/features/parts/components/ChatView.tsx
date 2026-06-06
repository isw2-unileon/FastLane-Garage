import React, { useState } from 'react';
import { Send, Bot, User, ArrowLeft, ShieldAlert } from 'lucide-react';

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
      text: '¡Hola! He recibido los componentes y servicios que has seleccionado en el configurador. El equipo técnico ya está revisando la viabilidad de la orden.',
      time: 'Justo ahora'
    },
    {
      id: 'msg-2',
      sender: 'bot',
      text: 'Soy el asistente automatizado del taller de Back-end. ¿Tienes alguna duda sobre los tiempos de instalación o compatibilidad de los accesorios eléctricos?',
      time: 'Justo ahora'
    }
  ]);
  const [inputValue, setInputValue] = useState('');

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim()) return;

    const userTime = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    
    // Mensaje del usuario
    const newUserMessage: Message = {
      id: `msg-user-${Date.now()}`,
      sender: 'user',
      text: inputValue,
      time: userTime
    };

    setMessages(prev => [...prev, newUserMessage]);
    setInputValue('');

    // Simulamos una respuesta estática del bot ya que aún no está conectado al Back-end
    setTimeout(() => {
      const botTime = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      setMessages(prev => [
        ...prev,
        {
          id: `msg-bot-${Date.now()}`,
          sender: 'bot',
          text: 'Entendido. Recibo tu mensaje, pero recuerda que actualmente me encuentro en modo simulación (Desconectado del Back-end). Almacenaré tu consulta para cuando la pasarela esté activa.',
          time: botTime
        }
      ]);
    }, 1000);
  };

  return (
    <div className="w-full h-[calc(100vh-72px)] bg-slate-950 flex flex-col">
      
      {/* Barra superior del chat */}
      <div className="px-6 py-3 bg-slate-900/40 backdrop-blur-md border-b border-slate-800 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button 
            onClick={onBackToGarage}
            className="p-1.5 rounded-lg bg-slate-950 border border-slate-800 text-slate-400 hover:text-cyan-400 hover:border-cyan-500/50 transition-all cursor-pointer"
            title="Volver al Garaje"
          >
            <ArrowLeft className="w-4 h-4" />
          </button>
          
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 bg-cyan-950 border border-cyan-500/30 rounded-full flex items-center justify-center relative">
              <Bot className="text-cyan-400 w-5 h-5" />
              {/* Indicador Online */}
              <span className="absolute bottom-0 right-0 w-2.5 h-2.5 bg-emerald-500 rounded-full border-2 border-slate-950 animate-pulse"></span>
            </div>
            <div>
              <h3 className="text-sm font-bold text-slate-200">Asistente Virtual Taller</h3>
              <p className="text-[10px] text-emerald-400 font-medium">Core Bot v1.0 • En Línea</p>
            </div>
          </div>
        </div>

        {/* Tag Informativo de integración futura */}
        <div className="hidden sm:flex items-center gap-2 bg-amber-950/20 border border-amber-900/40 text-amber-400 px-3 py-1 rounded-md text-[10px] uppercase font-mono tracking-wider">
          <ShieldAlert className="w-3.5 h-3.5" />
          Modo Maqueta / Offline
        </div>
      </div>

      {/* Cuerpo de Mensajes */}
      <div className="flex-1 overflow-y-auto p-6 flex flex-col gap-4 custom-scrollbar bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-slate-900/10 via-transparent to-transparent">
        {messages.map((msg) => {
          const isBot = msg.sender === 'bot';
          return (
            <div 
              key={msg.id} 
              className={`flex gap-3 max-w-[75%] ${isBot ? 'self-start' : 'self-end flex-row-reverse'}`}
            >
              {/* Avatar */}
              <div className={`w-7 h-7 rounded-full flex items-center justify-center flex-shrink-0 text-xs border ${
                isBot ? 'bg-slate-900 border-slate-800 text-cyan-400' : 'bg-cyan-950 border-cyan-800 text-cyan-300'
              }`}>
                {isBot ? <Bot className="w-3.5 h-3.5" /> : <User className="w-3.5 h-3.5" />}
              </div>

              {/* Contenido de la burbuja */}
              <div className="flex flex-col gap-1">
                <div className={`rounded-xl p-3 text-xs leading-relaxed border ${
                  isBot 
                    ? 'bg-slate-900/60 border-slate-800 text-slate-300 rounded-tl-none' 
                    : 'bg-gradient-to-br from-cyan-950/60 to-slate-900/60 border-cyan-900/50 text-cyan-100 rounded-tr-none'
                }`}>
                  {msg.text}
                </div>
                <span className={`text-[9px] text-slate-500 font-mono ${!isBot && 'text-right'}`}>
                  {msg.time}
                </span>
              </div>
            </div>
          );
        })}
      </div>

      {/* Formulario de Input inferior */}
      <form 
        onSubmit={handleSendMessage}
        className="p-4 bg-slate-900/20 border-t border-slate-900 flex gap-3 items-center"
      >
        <input 
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          placeholder="Escribe un mensaje al bot del taller..."
          className="flex-1 bg-slate-950 border border-slate-800 rounded-lg px-4 py-2.5 text-xs text-slate-300 placeholder-slate-600 focus:outline-none focus:border-cyan-500/50 transition-colors"
        />
        <button 
          type="submit"
          className="p-2.5 bg-cyan-950 text-cyan-400 border border-cyan-800 hover:bg-cyan-400 hover:text-slate-950 rounded-lg transition-all flex items-center justify-center cursor-pointer disabled:opacity-40"
          disabled={!inputValue.trim()}
        >
          <Send className="w-4 h-4" />
        </button>
      </form>

    </div>
  );
}