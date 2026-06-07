export interface ServiceItem {
  id: string;
  name: string;
  price: number;
  duration: string;
  description: string;
  category: 'wheels' | 'engine' | 'electricity' | 'brakes';
}

export interface VehicleSessionData {
  vehicle_brand: string;
  vehicle_model: string;
  vehicle_year: string;
}

export const GARAGE_SERVICES: ServiceItem[] = [
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

// --- AÑADE ESTO ABAJO ---
// Esto es lo que conecta el Frontend con tu Backend en Go
const API_URL = "/api";

export const chatApi = {
  createSession: async (vehicleData: VehicleSessionData) => {
    const response = await fetch(`${API_URL}/chat/sessions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(vehicleData),
    });
    if (!response.ok) throw new Error('Error al crear sesión');
    return response.json();
  },

  sendMessage: async (sessionId: number, text: string) => {
    const response = await fetch(`${API_URL}/chat/sessions/${sessionId}/messages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text }),
    });
    if (!response.ok) throw new Error('Error al enviar mensaje');
    return response.json();
  },

  getHistory: async (sessionId: number) => {
    const response = await fetch(`${API_URL}/chat/sessions/${sessionId}`);
    if (!response.ok) throw new Error('Error al obtener historial');
    return response.json();
  }
};
