
interface ServiceItem {
  id: string;
  name: string;
  duration: string;
  description: string;
}

const SERVICES_BY_ZONE: { ruedas: ServiceItem[]; iluminacion: ServiceItem[]; puertas: ServiceItem[] } = {
  ruedas: [
    { id: 'srv-w1', name: 'Cambio de Neumáticos (X4)', duration: '1h 30min', description: 'Montaje, desmontaje y equilibrado de los cuatro neumáticos de fábrica.' },
    { id: 'srv-w2', name: 'Alineación de Dirección 3D', duration: '45min', description: 'Ajuste de ángulos de convergencia para asegurar una pisada óptima.' },
    { id: 'srv-w3', name: 'Equilibrado de Ejes', duration: '30min', description: 'Corrección de contrapesos para eliminar vibraciones en el volante.' }
  ],
  iluminacion: [
    { id: 'srv-i1', name: 'Pulido y Lacado de Ópticas', duration: '1h', description: 'Eliminación del desgaste por UV para recuperar la máxima claridad de los faros.' },
    { id: 'srv-i2', name: 'Conversión a Sistema LED Pro', duration: '30min', description: 'Sustitución de bombillas estándar por tecnología LED de alta intensidad.' },
    { id: 'srv-i3', name: 'Regulación de Haz de Luz', duration: '15min', description: 'Calibración de altura de faros para evitar deslumbramientos y pasar ITV.' }
  ],
  puertas: [
    { id: 'srv-p1', name: 'Ajuste y Engrase de Bisagras', duration: '20min', description: 'Eliminación de ruidos y alineación del cierre de las puertas laterales.' },
    { id: 'srv-p2', name: 'Insonorización Acústica Premium', duration: '2h 30min', description: 'Instalación de planchas aislantes dentro del panel para reducir ruido de rodadura.' },
    { id: 'srv-p3', name: 'Reparación de Elevalunas Eléctrico', duration: '1h', description: 'Sustitución del motor o guías del sistema de subida del cristal.' }
  ]
};

interface DashboardProps {
  selectedZone: string | null;
  onAddToCart?: (id: string, name: string) => void;
  onRemoveFromCart?: (id: string) => void;
  cartItems?: Array<{ id: string }>;
}

export default function Dashboard({ selectedZone, onAddToCart, onRemoveFromCart, cartItems = [] }: DashboardProps) {
  
  const getCurrentServices = (): ServiceItem[] => {
    if (selectedZone === 'ruedas') return SERVICES_BY_ZONE.ruedas;
    if (selectedZone === 'iluminacion') return SERVICES_BY_ZONE.iluminacion;
    if (selectedZone === 'puertas') return SERVICES_BY_ZONE.puertas;
    return [...SERVICES_BY_ZONE.ruedas, ...SERVICES_BY_ZONE.iluminacion, ...SERVICES_BY_ZONE.puertas];
  };

  const currentServices = getCurrentServices();

  const getZoneTitle = (zone: string | null) => {
    switch (zone) {
      case 'ruedas': return 'Eje & Neumáticos';
      case 'iluminacion': return 'Sistema de Iluminación';
      case 'puertas': return 'Carrocería & Accesos (Puertas)';
      default: return 'Todos los Servicios del Taller';
    }
  };

  return (
    
    <div className="w-full bg-slate-900/40 backdrop-blur-md border-t border-slate-800 p-6 flex flex-col gap-4">
      
      {/* Cabecera del catálogo */}
      <div className="flex justify-between items-center border-b border-slate-800 pb-3">
        <div>
          <h2 className="text-cyan-400 font-bold tracking-wider text-sm uppercase">
            {getZoneTitle(selectedZone)}
          </h2>
          <p className="text-xs text-slate-400">
            {selectedZone ? 'Operaciones mecánicas calibradas para el área seleccionada' : 'Catálogo completo de intervenciones disponibles'}
          </p>
        </div>
        <span className="text-xs bg-cyan-950 text-cyan-400 border border-cyan-800 px-3 py-1 rounded">
          {currentServices.length} Opciones
        </span>
      </div>

     
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {currentServices.map((service) => {
          const isAdded = cartItems.some(item => item.id === service.id);
          return (
            <div 
              key={service.id}
              className="w-full h-32 bg-slate-950 border border-slate-800 rounded p-3 flex flex-col justify-between hover:border-cyan-500/50 transition-colors group"
            >
              <div>
                <div className="flex justify-between items-start gap-2">
                  <h3 className="text-slate-200 font-medium text-sm line-clamp-1 group-hover:text-cyan-400 transition-colors">
                    {service.name}
                  </h3>
                  <span className="text-[10px] text-slate-500 bg-slate-900 px-1.5 py-0.5 rounded border border-slate-800 font-mono flex-shrink-0">
                    {service.duration}
                  </span>
                </div>
                <p className="text-xs text-slate-400 mt-1 line-clamp-2">
                  {service.description}
                </p>
              </div>

              {/* Acción inferior */}
              <div className="flex justify-end items-center mt-2 pt-2 border-t border-slate-900">
                {isAdded ? (
                  <button 
                    onClick={() => onRemoveFromCart && onRemoveFromCart(service.id)}
                    className="text-[10px] uppercase tracking-wider px-3 py-1 rounded font-bold bg-red-950 text-red-400 border border-red-900 hover:bg-red-500 hover:text-white transition-all"
                  >
                    Quitar
                  </button>
                ) : (
                  <button 
                    onClick={() => onAddToCart && onAddToCart(service.id, service.name)}
                    className="text-[10px] uppercase tracking-wider px-3 py-1 rounded font-bold bg-cyan-950 text-cyan-400 border border-cyan-800 hover:bg-cyan-400 hover:text-slate-950 transition-all"
                  >
                    Contratar
                  </button>
                )}
              </div>
            </div>
          );
        })}
      </div>

    </div>
  );
}