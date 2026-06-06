export interface ServiceItem {
  id: string;
  name: string;
  price: number;
  duration: string; // Ej: "1.5 horas"
  description: string;
  category: 'wheels' | 'engine' | 'electricity' | 'brakes';
}

export const GARAGE_SERVICES: ServiceItem[] = [
  // Categoría Ruedas (Sustituye a los recambios de ruedas)
  {
    id: 'srv-1',
    name: 'Cambio de Neumáticos (X4)',
    price: 80,
    duration: '1h 30min',
    description: 'Montaje, desmontaje y equilibrado de los cuatro neumáticos.',
    category: 'wheels'
  },
  {
    id: 'srv-2',
    name: 'Alineación de Dirección 3D',
    price: 55,
    duration: '45min',
    description: 'Ajuste de ángulos de las ruedas para asegurar la pisada correcta y evitar desgaste.',
    category: 'wheels'
  },
  {
    id: 'srv-3',
    name: 'Equilibrado de Ejes',
    price: 30,
    duration: '30min',
    description: 'Eliminación de vibraciones en el volante a altas velocidades.',
    category: 'wheels'
  }
];