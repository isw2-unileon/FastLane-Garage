import React from 'react';
import { ShoppingCart, LayoutGrid, BarChart3 } from 'lucide-react';

interface HeaderProps {
  cartCount: number;
  onToggleCart: () => void;
  currentView: 'garage' | 'analytics';
  setView: (view: 'garage' | 'analytics') => void;
}

export default function Header({ cartCount, onToggleCart, currentView, setView }: HeaderProps) {
  return (
    <header className="fixed top-0 left-0 w-full z-50 flex items-center justify-between px-6 py-4 bg-gradient-to-b from-[#070a13] to-transparent">
      {/* Izquierda: Logo */}
      <div className="flex items-center gap-3 cursor-pointer" onClick={() => setView('garage')}>
        <div className="w-8 h-8 bg-cyan-500 rounded-lg flex items-center justify-center shadow-[0_0_15px_rgba(6,182,212,0.5)]">
          <LayoutGrid className="text-black w-5 h-5" />
        </div>
        <h1 className="text-white font-bold tracking-[0.2em] text-lg uppercase">Garage App</h1>
      </div>

      {/* Centro: Línea de espectro (Decorativa) */}
      <div className="hidden md:flex flex-1 max-w-md mx-10 items-center justify-center opacity-50">
        <div className="w-full h-[2px] bg-gradient-to-r from-transparent via-cyan-500 to-transparent relative">
          <div className="absolute -top-1 left-1/2 -translate-x-1/2 w-24 h-3 bg-cyan-500/20 blur-md rounded-full"></div>
        </div>
      </div>

      {/* Derecha: Acciones filtradas (Solo Estadísticas y Carrito) */}
      <div className="flex items-center gap-4">
        
        {/* 📊 BOTÓN DE ESTADÍSTICAS */}
        <button 
          onClick={() => setView(currentView === 'analytics' ? 'garage' : 'analytics')}
          title="Ver estadísticas"
          className={`p-2 rounded-full border transition-all cursor-pointer focus:outline-none ${
            currentView === 'analytics'
              ? 'bg-cyan-500 text-slate-950 border-cyan-400 shadow-[0_0_15px_rgba(6,182,212,0.4)]'
              : 'bg-slate-900/60 text-slate-400 border-slate-800 hover:border-cyan-500/50 hover:text-cyan-400'
          }`}
        >
          <BarChart3 className="w-4 h-4" />
        </button>

        {/* 🛒 BOTÓN INTERACTIVO DEL CARRITO */}
        <button 
          onClick={onToggleCart}
          className="flex items-center gap-2 bg-cyan-950/40 border border-cyan-500/30 px-3 py-1.5 rounded-full cursor-pointer hover:bg-cyan-900/60 hover:border-cyan-400/60 transition-all focus:outline-none"
        >
          <ShoppingCart className="text-cyan-400 w-4 h-4" />
          <span className="text-cyan-400 text-xs font-bold font-mono">{cartCount} Items</span>
        </button>

      </div>
    </header>
  );
}