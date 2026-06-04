import React, { useState } from 'react';
import Header from './components/layout/Header';
import Scene3D from './features/garage/components/Scene3D';
import Dashboard from './features/parts/components/Dashboard';
import AnalyticsView from './features/parts/components/AnalyticsView';
import ChatView from './features/parts/components/ChatView'; // 👈 Importamos la nueva vista del chat

const ELECTRICAL_ACCESSORIES = [
  { id: 'elec-1', name: 'Módulo Ambient Light RGB' },
  { id: 'elec-2', name: 'Dashcam Frontal 4K Pro' },
  { id: 'elec-3', name: 'Cargador Inalámbrico QI MagSafe' },
  { id: 'elec-4', name: 'Pantalla Infotainment 10.2"' },
  { id: 'elec-6', name: 'Sensor de Aparcamiento Ultrasonido' },
  { id: 'elec-7', name: 'Inversor de Corriente 220V' }
];

interface CartItem {
  id: string;
  name: string;
}

export default function App() {
  const [selectedZone, setSelectedZone] = useState<string | null>(null);
  const [cart, setCart] = useState<CartItem[]>([]);
  const [isCartOpen, setIsCartOpen] = useState(false);
  
  // 🧭 CONTROL DE RUTAS AMPLIADO: 'garage' | 'analytics' | 'chat'
  const [view, setView] = useState<'garage' | 'analytics' | 'chat'>('garage');

  const addToCart = (id: string, name: string) => {
    if (cart.some(item => item.id === id)) return;
    setCart([...cart, { id, name }]);
  };

  const removeFromCart = (id: string) => {
    setCart(cart.filter(item => item.id !== id));
  };

  // 🚀 LOGICA DE ENVIO AL ASESOR: Ahora limpia el carro y salta a la pantalla de Chat
  const handleSendToAdvisor = () => {
    setCart([]);
    setIsCartOpen(false);
    setView('chat'); // Redirección automática a la vista del Bot
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white flex flex-col overflow-y-auto overflow-x-hidden font-sans select-none relative">
      
      {/* Header unificado (Si estamos en el chat, marcamos analíticas como apagado) */}
      <Header 
        cartCount={cart.length} 
        onToggleCart={() => setIsCartOpen(!isCartOpen)} 
        currentView={view === 'analytics' ? 'analytics' : 'garage'}
        setView={(targetView) => setView(targetView)}
      />

      {/* DESPLEGABLE FLOTANTE DEL CARRITO */}
      {isCartOpen && (
        <div className="absolute top-[68px] right-6 w-80 bg-slate-900/95 backdrop-blur-xl border border-slate-800 rounded-xl p-4 shadow-2xl z-50 flex flex-col max-h-[450px]">
          <div className="flex justify-between items-center border-b border-slate-800 pb-2 mb-3">
            <h3 className="font-bold text-sm tracking-wider text-cyan-400 uppercase">Resumen del Pedido</h3>
            <span className="text-[10px] text-slate-500 font-mono">{cart.length} items</span>
          </div>

          <div className="flex-1 overflow-y-auto pr-1 flex flex-col gap-2 custom-scrollbar">
            {cart.length === 0 ? (
              <p className="text-xs text-slate-500 text-center py-6">El carrito está vacío.</p>
            ) : (
              cart.map((item) => (
                <div key={item.id} className="bg-slate-950/60 border border-slate-800/60 rounded p-2.5 flex justify-between items-center group">
                  <span className="text-xs text-slate-300 truncate font-medium max-w-[80%]">{item.name}</span>
                  <button onClick={() => removeFromCart(item.id)} className="text-[10px] text-slate-500 hover:text-red-400 font-medium px-1.5 py-0.5 rounded transition-colors">
                    Quitar
                  </button>
                </div>
              ))
            )}
          </div>

          {cart.length > 0 && (
            <div className="border-t border-slate-800 pt-3 mt-3">
              <button onClick={handleSendToAdvisor} className="w-full text-center text-xs py-2 bg-gradient-to-r from-cyan-500 to-cyan-600 hover:from-cyan-400 hover:to-cyan-500 text-slate-950 font-bold uppercase tracking-wider rounded-lg transition-all shadow-[0_0_15px_rgba(6,182,212,0.3)]">
                Enviar al asesor
              </button>
            </div>
          )}
        </div>
      )}

      {/* ENRUTADOR PRINCIPAL DE VISTAS */}
      <div className="flex-1 flex flex-col pt-[72px]">
        
        {view === 'analytics' && (
          <AnalyticsView />
        )}

        {view === 'chat' && (
          /* 💬 RUTA DEL CHAT CON BOTÓN INTERNO PARA REGRESAR AL GARAJE */
          <ChatView onBackToGarage={() => setView('garage')} />
        )}

        {view === 'garage' && (
          <>
            {/* BLOQUE SUPERIOR FIJO */}
            <div className="h-[58vh] w-full flex flex-row border-b border-slate-900 flex-shrink-0">
              <div className="flex-1 h-full relative">
                <Scene3D selectedZone={selectedZone} onSelectZone={(zone) => setSelectedZone(zone)} />
              </div>

              <div className="w-80 h-full bg-slate-900/50 backdrop-blur-md border-l border-slate-800 p-4 flex flex-col flex-shrink-0">
                <div className="mb-3">
                  <h2 className="text-cyan-400 font-bold tracking-wider text-xs uppercase">Componentes Eléctricos</h2>
                  <p className="text-[10px] text-slate-400">Accesorios y actualizaciones de sistema</p>
                </div>

                <div className="flex-1 overflow-y-auto pr-1 flex flex-col gap-2 custom-scrollbar">
                  {ELECTRICAL_ACCESSORIES.map((item) => {
                    const isAdded = cart.some(c => c.id === item.id);
                    return (
                      <div key={item.id} className="bg-slate-950 border border-slate-800 rounded p-2.5 flex justify-between items-center hover:border-cyan-500/30 transition-colors group">
                        <div className="flex flex-col max-w-[65%]">
                          <span className="text-xs text-slate-300 group-hover:text-cyan-400 transition-colors font-medium truncate">{item.name}</span>
                          <span className="text-[10px] text-slate-500">Categoría: Electrónica</span>
                        </div>
                        <div className="flex items-center gap-1.5">
                          {isAdded ? (
                            <button onClick={() => removeFromCart(item.id)} className="text-[9px] uppercase tracking-wider px-2 py-0.5 border border-red-900/40 bg-red-950/20 text-red-400 rounded font-medium hover:bg-red-500 hover:text-white transition-colors">Quitar</button>
                          ) : (
                            <button onClick={() => addToCart(item.id, item.name)} className="text-[9px] uppercase tracking-widest px-2 py-0.5 border border-cyan-900/30 bg-cyan-950/50 text-cyan-400 rounded font-medium hover:bg-cyan-400 hover:text-slate-950 transition-all">+ Añadir</button>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>

            {/* DASHBOARD INFERIOR */}
            <div className="w-full flex-shrink-0">
              <Dashboard selectedZone={selectedZone} onAddToCart={addToCart} onRemoveFromCart={removeFromCart} cartItems={cart} />
            </div>
          </>
        )}

      </div>
    </div>
  );
}