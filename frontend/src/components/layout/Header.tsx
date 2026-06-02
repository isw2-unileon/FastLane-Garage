import { ShoppingCart, Settings, User, LogOut, LayoutGrid } from 'lucide-react';

export default function Header() {
  return (
    <header className="fixed top-0 left-0 w-full z-50 flex items-center justify-between px-6 py-4 bg-gradient-to-b from-[#070a13] to-transparent">
      {/* Izquierda: Logo */}
      <div className="flex items-center gap-3">
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

      {/* Derecha: Acciones */}
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2 bg-cyan-950/40 border border-cyan-500/30 px-3 py-1.5 rounded-full cursor-pointer hover:bg-cyan-900/60 transition-all">
          <ShoppingCart className="text-cyan-400 w-4 h-4" />
          <span className="text-cyan-400 text-xs font-bold font-mono">3 Items</span>
        </div>
        <div className="flex items-center gap-4 text-slate-400">
          <Settings className="w-5 h-5 hover:text-white cursor-pointer transition-colors" />
          <User className="w-5 h-5 hover:text-white cursor-pointer transition-colors" />
          <LogOut className="w-5 h-5 hover:text-red-400 cursor-pointer transition-colors" />
        </div>
      </div>
    </header>
  );
}