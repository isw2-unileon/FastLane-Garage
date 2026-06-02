import { ShoppingCart } from 'lucide-react';

interface DashboardProps {
  selectedZone: string | null;
}

export default function Dashboard({ selectedZone }: DashboardProps) {
  const zonaLimpia = selectedZone ? selectedZone.replace('_Inter', '').replace(/_/g, ' ') : "";

  const llantas = [
    { id: 1, nombre: 'Llantas Aleación 19" Aero', precio: '349 €', img: 'https://images.unsplash.com/photo-1594731826601-57497ec09099?w=150&h=150&fit=crop' },
    { id: 2, nombre: 'Llantas Deportivas 20" Black', precio: '415 €', img: 'https://images.unsplash.com/photo-1551009175-15bdf9dcb580?w=150&h=150&fit=crop' },
    { id: 3, nombre: 'Llantas Forjadas 19" Bronze', precio: '510 €', img: 'https://images.unsplash.com/photo-1580273916550-e323be2ae537?w=150&h=150&fit=crop' },
  ];

  const neumaticos = [
    { id: 4, nombre: 'Michelin Sport 4S', precio: '210 €/u', img: 'https://images.unsplash.com/photo-1549411380-496677f52554?w=150&h=150&fit=crop' },
    { id: 5, fontCustom: true, nombre: 'Pirelli P Zero', precio: '195 €/u', img: 'https://images.unsplash.com/photo-1620311210427-463e26f58273?w=150&h=150&fit=crop' },
    { id: 6, nombre: 'Continental SportContact 6', precio: '205 €/u', img: 'https://images.unsplash.com/photo-1532634922-8fe0b757fb13?w=150&h=150&fit=crop' },
  ];

  return (
    <div className="w-full h-full bg-[#070a13]/90 border-t border-cyan-500/20 backdrop-blur-md px-8 py-6 flex flex-col justify-between">
      {selectedZone ? (
        <div className="h-full flex flex-col justify-between">
          
          {/* Encabezado del HUD */}
          <div className="flex items-center gap-3 mb-4 text-xs font-bold tracking-[0.25em] text-slate-400 font-mono">
            <div className="w-1.5 h-4 bg-cyan-500 shadow-[0_0_10px_rgba(6,182,212,0.7)]"></div>
            <span className="uppercase">TREN DE RODAJE –</span>
            <span className="text-cyan-400 uppercase font-sans tracking-[0.2em]">{zonaLimpia}</span>
          </div>

          {/* Contenedor de las dos columnas con scroll oculto */}
          <div className="flex-1 flex gap-10 overflow-x-auto overflow-y-hidden pb-2 style-scrollbar-hidden">
            
            {/* COLUMNA: LLANTAS */}
            <div className="flex flex-col gap-3">
              <div className="text-[10px] font-bold text-slate-500 tracking-widest flex items-center gap-2 font-mono">
                <span className="w-1 h-1 rounded-full bg-cyan-500/60"></span> LLANTAS
              </div>
              <div className="flex gap-4">
                {llantas.map((item) => (
                  <PartCard key={item.id} item={item} />
                ))}
              </div>
            </div>

            {/* Separador vertical estético */}
            <div className="w-[1px] h-40 bg-gradient-to-b from-transparent via-slate-800 to-transparent self-center"></div>

            {/* COLUMNA: NEUMÁTICOS */}
            <div className="flex flex-col gap-3">
              <div className="text-[10px] font-bold text-slate-500 tracking-widest flex items-center gap-2 font-mono">
                <span className="w-1 h-1 rounded-full bg-cyan-500/60"></span> NEUMÁTICOS
              </div>
              <div className="flex gap-4">
                {neumaticos.map((item) => (
                  <PartCard key={item.id} item={item} isSelectedStyle={item.fontCustom} />
                ))}
              </div>
            </div>

          </div>
        </div>
      ) : (
        <div className="w-full h-full flex flex-col items-center justify-center">
          <p className="text-cyan-500/40 text-[11px] tracking-[0.5em] uppercase font-mono animate-pulse">
            ▲ SELECCIONA UNA ZONA ACTIVA EN EL VISOR HOLOGRÁFICO ▲
          </p>
        </div>
      )}
    </div>
  );
}

function PartCard({ item, isSelectedStyle = false }: { item: any; isSelectedStyle?: boolean }) {
  return (
    <div className={`w-48 h-[175px] bg-[#0b1220]/70 border rounded-xl p-3 flex flex-col justify-between transition-all duration-300 group
      ${isSelectedStyle 
        ? 'border-cyan-400/80 shadow-[0_0_15px_rgba(6,182,212,0.1)] bg-[#0f1d35]/80' 
        : 'border-slate-800/80 hover:border-cyan-500/30 hover:bg-[#0e1729]'
      }`}
    >
      <div className="w-full h-20 rounded-lg bg-[#05080f] border border-slate-900 overflow-hidden flex items-center justify-center relative">
        <img 
          src={item.img} 
          alt={item.nombre} 
          className="w-full h-full object-cover opacity-75 group-hover:scale-105 group-hover:opacity-100 transition-all duration-300 mix-blend-screen"
        />
      </div>

      <div className="mt-1 px-0.5">
        <h4 className="text-slate-300 text-[11px] font-medium tracking-tight truncate group-hover:text-white transition-colors">
          {item.nombre}
        </h4>
        <p className="text-cyan-400 text-[10px] font-mono font-bold mt-0.5">
          {item.precio}
        </p>
      </div>

      <button className={`w-full py-1 border rounded-lg text-[10px] font-bold uppercase tracking-wider flex items-center justify-center gap-2 transition-all duration-200
        ${isSelectedStyle
          ? 'bg-cyan-500 text-black border-transparent font-extrabold shadow-[0_0_10px_rgba(6,182,212,0.3)]'
          : 'bg-[#0a101d] border-slate-800 text-slate-400 hover:border-cyan-500/40 hover:text-cyan-400 hover:bg-cyan-950/20'
        }`}
      >
        <ShoppingCart className="w-3 h-3" />
        Añadir
      </button>
    </div>
  );
}