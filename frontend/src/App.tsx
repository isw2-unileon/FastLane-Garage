import { useState, useEffect } from 'react';
import Scene3D from './features/garage/components/Scene3D';
import Dashboard from './features/parts/components/Dashboard';
import Header from './components/layout/Header';

interface PartItem {
  id: string;
  name: string;
  price: number;
  image: string;       
  description?: string; 
}

export default function App() {
  const [selectedZone, setSelectedZone] = useState<string | null>(null);
  const [electricParts, setElectricParts] = useState<PartItem[]>([]);
  const [selectedOil, setSelectedOil] = useState<string>('5w30');

  useEffect(() => {
    const mockElectricParts: PartItem[] = [
      { 
        id: 'elec_neon', 
        name: 'Neon Inferior RGB', 
        price: 45, 
        image: 'https://images.unsplash.com/photo-1552519507-da3b142c6e3d?w=150&q=80', 
        description: 'Kit LED sincronizado' 
      },
      { 
        id: 'elec_bat', 
        name: 'Batería de Litio', 
        price: 180, 
        image: 'https://images.unsplash.com/photo-1620214948434-5515e065bc54?w=150&q=80',
        description: 'Gel de alta capacidad' 
      },
      { 
        id: 'elec_ecu', 
        name: 'Reprogramación ECU', 
        price: 320, 
        image: 'https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=150&q=80',
        description: 'Optimización Stage 1' 
      }
    ];
    setElectricParts(mockElectricParts);
  }, []);

  const handleAddToCart = (item: { id: string; name: string; price: number }) => {
    console.log(`Añadiendo al carrito: ${item.name} (${item.price}€)`);
    alert(`Añadido al carrito: ${item.name}`);
  };

  return (
    <div className="w-screen h-screen flex flex-col bg-[#030712] overflow-hidden select-none antialiased text-slate-200">
      {/* Barra superior fija */}
      <Header />

      {/* CONTENEDOR FLEX DE LA ZONA SUPERIOR */}
      <div className="w-full h-[58vh] flex mt-14 border-b border-slate-900">
        
        {/* 🌌 SUB-SECCIÓN IZQUIERDA: Visor 3D */}
        <div className="w-3/4 h-full relative">
          <Scene3D selectedZone={selectedZone} onSelectZone={setSelectedZone} />
          
          {/* Decoración HUD */}
          <div className="absolute top-1/2 left-6 -translate-y-1/2 flex flex-col gap-6 opacity-20 pointer-events-none">
             {[...Array(4)].map((_, i) => (
               <div key={i} className="w-8 h-[1px] bg-cyan-500"></div>
             ))}
          </div>
        </div>

        {/* 🎛️ SUB-SECCIÓN DERECHA: Panel Técnico Configurable */}
        <div className="w-1/4 h-full bg-[#050b18]/95 border-l border-slate-850 p-4 flex flex-col gap-4 backdrop-blur-md overflow-y-auto z-20">

          {/* ⚡ COMPLEMENTOS ELÉCTRICOS (Scroll Vertical con imagen a la izquierda) */}
          <div className="flex flex-col gap-2">
            <div className="flex items-center gap-2 text-cyan-400">
              <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              <h3 className="text-xs font-semibold tracking-wider uppercase">Complementos Eléctricos</h3>
            </div>
            
            {/* 📜 Contenedor con Scroll Vertical acotado */}
            <div className="flex flex-col gap-2 max-h-[200px] overflow-y-auto pr-1 scrollbar-thin scrollbar-thumb-slate-800">
              {electricParts.map((part) => (
                <div 
                  key={part.id} 
                  className="w-full bg-slate-900/40 border border-slate-800 p-2 rounded flex gap-3 items-center"
                >
                  {/* 📸 Izquierda: Imagen del producto fija */}
                  <div className="w-16 h-16 bg-slate-950 rounded overflow-hidden flex-shrink-0 border border-slate-850">
                    <img 
                      src={part.image} 
                      alt={part.name} 
                      className="w-full h-full object-cover opacity-80"
                    />
                  </div>

                  {/* ➡️ Derecha: Nombre, descripción, precio y botón alineados */}
                  <div className="flex-1 flex flex-col justify-between h-16 min-w-0">
                    <div className="flex justify-between items-start gap-1">
                      <div className="min-w-0">
                        <div className="text-[11px] font-bold tracking-wide truncate text-slate-200">{part.name}</div>
                        <div className="text-[9px] text-slate-400 truncate">{part.description}</div>
                      </div>
                      <span className="text-[10px] font-mono font-bold text-cyan-400 flex-shrink-0">
                        {part.price}€
                      </span>
                    </div>

                    {/* Botón Carrito compacto */}
                    <button 
                      onClick={() => handleAddToCart({ id: part.id, name: part.name, price: part.price })}
                      className="w-full bg-cyan-950/40 hover:bg-cyan-500 border border-cyan-500/30 hover:border-cyan-400 text-cyan-400 hover:text-slate-950 text-[9px] font-bold py-1 rounded transition-all flex items-center justify-center gap-1"
                    >
                      <svg className="w-2.5 h-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                      </svg>
                      Añadir al carrito
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>


          {/* 🛢️ CAMBIO DE ACEITE */}
          <div className="flex flex-col gap-2">
            <div className="flex items-center gap-2 text-amber-500">
              <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
              </svg>
              <h3 className="text-xs font-semibold tracking-wider uppercase">Cambio de Aceite</h3>
            </div>

            <div className="bg-slate-900/20 border border-slate-800/60 rounded p-2.5 flex flex-col gap-2.5">
              <div>
                <div className="flex justify-between text-[10px] mb-0.5 text-slate-400">
                  <span>Vida útil del aceite</span>
                  <span className="text-amber-500 font-mono font-bold">12%</span>
                </div>
                <div className="w-full bg-slate-950 h-1 rounded-full overflow-hidden border border-slate-800">
                  <div className="bg-gradient-to-r from-red-500 to-amber-500 h-full w-[12%]" />
                </div>
              </div>

              <div className="flex flex-col gap-1">
                <label className="flex items-center justify-between bg-slate-950/40 p-1.5 rounded border border-slate-800/60 cursor-pointer hover:border-amber-500/30 transition-all text-[10px]">
                  <div className="flex items-center gap-1.5">
                    <input 
                      type="radio" 
                      name="oilType" 
                      checked={selectedOil === '5w30'} 
                      onChange={() => setSelectedOil('5w30')}
                      className="accent-amber-500" 
                    />
                    <span className="text-slate-300">Sintético Premium 5W-30</span>
                  </div>
                  <span className="font-mono text-amber-500 font-bold">75€</span>
                </label>
                <label className="flex items-center justify-between bg-slate-950/40 p-1.5 rounded border border-slate-800/60 cursor-pointer hover:border-amber-500/30 transition-all text-[10px]">
                  <div className="flex items-center gap-1.5">
                    <input 
                      type="radio" 
                      name="oilType" 
                      checked={selectedOil === 'competicion'} 
                      onChange={() => setSelectedOil('competicion')}
                      className="accent-amber-500" 
                    />
                    <span className="text-slate-300">Competición Alta Viscosidad</span>
                  </div>
                  <span className="font-mono text-amber-500 font-bold">120€</span>
                </label>
              </div>

              <button 
                onClick={() => handleAddToCart({
                  id: `oil_${selectedOil}`, 
                  name: `Cambio de Aceite (${selectedOil === '5w30' ? '5W-30' : 'Competición'})`, 
                  price: selectedOil === '5w30' ? 75 : 120
                })}
                className="w-full bg-gradient-to-r from-amber-600 to-amber-500 hover:from-amber-500 hover:to-amber-400 text-slate-950 text-[11px] font-bold py-1.5 px-3 rounded transition-all shadow-md shadow-amber-950/10"
              >
                Añadir Cambio a la Lista
              </button>
            </div>
          </div>

        </div>
      </div>

      {/* ZONA INFERIOR: Catálogo dinámico */}
      <div className="w-full h-[42vh] relative z-10">
        <Dashboard selectedZone={selectedZone} />
      </div>
    </div>
  );
}