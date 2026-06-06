import { TrendingUp, Clock, Package, CheckCircle2 } from 'lucide-react';

const MOST_REQUESTED_SERVICES = [
  { name: 'Cambio de Neumáticos (X4)', category: 'Mecánica', count: 142, percentage: 92, color: 'from-cyan-500 to-blue-600' },
  { name: 'Conversión a Sistema LED Pro', category: 'Iluminación', count: 118, percentage: 78, color: 'from-purple-500 to-indigo-600' },
  { name: 'Alineación de Dirección 3D', category: 'Mecánica', count: 95, percentage: 64, color: 'from-emerald-500 to-teal-600' },
  { name: 'Módulo Ambient Light RGB', category: 'Electrónica', count: 84, percentage: 55, color: 'from-pink-500 to-rose-600' },
  { name: 'Insonorización Acústica Premium', category: 'Carrocería', count: 41, percentage: 28, color: 'from-amber-500 to-orange-600' },
];

export default function AnalyticsView() {
  return (
    <div className="w-full min-h-[calc(100vh-72px)] bg-slate-950 p-6 flex flex-col gap-6">
      
      {/* Encabezado */}
      <div>
        <h2 className="text-cyan-400 font-bold tracking-wider text-sm uppercase">Panel de Control & Analíticas</h2>
        <p className="text-xs text-slate-400">Métricas de rendimiento y demandas de servicios del taller en tiempo real</p>
      </div>

      {/* Tarjetas de Resumen Rápido (KPIs) */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {[
          { title: 'Total Solicitudes', value: '480', sub: '+12% este mes', icon: CheckCircle2, color: 'text-cyan-400' },
          { title: 'Servicio Estrella', value: 'Neumáticos X4', sub: '92% de demanda', icon: TrendingUp, color: 'text-purple-400' },
          { title: 'Tiempo Medio Caja', value: '54 min', sub: '-8m de optimización', icon: Clock, color: 'text-emerald-400' },
          { title: 'Accesorios Instalados', value: '125 u.', sub: 'Stock estable', icon: Package, color: 'text-amber-400' },
        ].map((kpi, idx) => (
          <div key={idx} className="bg-slate-900/40 backdrop-blur-md border border-slate-800 rounded-xl p-4 flex justify-between items-center">
            <div>
              <span className="text-xs text-slate-500 uppercase font-semibold tracking-wider">{kpi.title}</span>
              <h4 className="text-2xl font-mono font-bold text-white mt-1">{kpi.value}</h4>
              <span className="text-[10px] text-slate-400 mt-1 block">{kpi.sub}</span>
            </div>
            <kpi.icon className={`w-8 h-8 opacity-70 ${kpi.color}`} />
          </div>
        ))}
      </div>

      {/* Sección Central de Gráficos */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        {/* Gráfico Principal: Servicios Más Pedidos */}
        <div className="lg:col-span-2 bg-slate-900/20 backdrop-blur-md border border-slate-800 rounded-xl p-5 flex flex-col gap-4">
          <div className="border-b border-slate-800 pb-2">
            <h3 className="text-xs font-bold uppercase tracking-wider text-slate-300">Ranking de Intervenciones Solicitadas</h3>
          </div>

          <div className="flex flex-col gap-4 py-2">
            {MOST_REQUESTED_SERVICES.map((srv, idx) => (
              <div key={idx} className="flex flex-col gap-1">
                <div className="flex justify-between items-end text-xs">
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-slate-600 text-[10px]">0{idx + 1}</span>
                    <span className="text-slate-300 font-medium">{srv.name}</span>
                    <span className="text-[9px] bg-slate-900 border border-slate-800 text-slate-400 px-1.5 py-0.2 rounded-full">{srv.category}</span>
                  </div>
                  <span className="font-mono font-bold text-cyan-400">{srv.count} órdenes</span>
                </div>
                {/* Contenedor barra */}
                <div className="w-full h-3 bg-slate-950 rounded-full overflow-hidden border border-slate-900 relative">
                  <div 
                    className={`h-full rounded-full bg-gradient-to-r ${srv.color} shadow-[0_0_10px_rgba(6,182,212,0.2)]`} 
                    style={{ width: `${srv.percentage}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Gráfico Secundario: Distribución por Categorías */}
        <div className="bg-slate-900/20 backdrop-blur-md border border-slate-800 rounded-xl p-5 flex flex-col justify-between">
          <div className="border-b border-slate-800 pb-2 mb-3">
            <h3 className="text-xs font-bold uppercase tracking-wider text-slate-300">Cuota de Rendimiento</h3>
          </div>
          
          <div className="flex-1 flex flex-col justify-center gap-4">
            {[
              { label: 'Mecánica Completa', val: '55%', width: 'w-[55%]', color: 'bg-cyan-500' },
              { label: 'Sistemas Eléctricos', val: '25%', width: 'w-[25%]', color: 'bg-purple-500' },
              { label: 'Carrocería & Cierres', val: '12%', width: 'w-[12%]', color: 'bg-emerald-500' },
              { label: 'Otros Ensayos', val: '8%', width: 'w-[8%]', color: 'bg-amber-500' },
            ].map((cat, idx) => (
              <div key={idx} className="flex flex-col gap-1.5">
                <div className="flex justify-between text-xs font-medium">
                  <span className="text-slate-400">{cat.label}</span>
                  <span className="text-white font-mono">{cat.val}</span>
                </div>
                <div className="w-full h-1.5 bg-slate-950 rounded-full">
                  <div className={`h-full rounded-full ${cat.color}`} style={{ width: cat.val }} />
                </div>
              </div>
            ))}
          </div>
        </div>

      </div>

    </div>
  );
}