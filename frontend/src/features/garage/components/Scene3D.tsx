import { Canvas } from '@react-three/fiber';
import { OrbitControls, Center } from '@react-three/drei';
import { Suspense } from 'react';
import CarModel from './CarModel';
import MapaModel from './MapaModel';

interface Scene3DProps {
  selectedZone: string | null;
  onSelectZone: (zone: string | null) => void;
}

export default function Scene3D({ selectedZone, onSelectZone }: Scene3DProps) {
  return (
    <div className="w-full h-full bg-[#020617] relative">
      {/* 🚀 CAMBIO 1: Bajamos la altura inicial (Y=2.2) y ajustamos la distancia (Z=12) */}
      <Canvas camera={{ position: [0, 10, 30], fov: 45 }}>
        {/* Iluminación general potente */}
        <ambientLight intensity={0.7} />
        <directionalLight position={[10, 15, 10]} intensity={1.5} castShadow />
        <directionalLight position={[-10, 10, -10]} intensity={0.6} />

        <Suspense fallback={null}>
          {/* Mantenemos tu Center que calcula bien los límites de los modelos */}
          <Center>
            <CarModel selectedZone={selectedZone} onSelectZone={onSelectZone} />
            <MapaModel />
          </Center>
        </Suspense>

        {/* 🛠️ AJUSTES CRÍTICOS DE CONTROL Y LIMITACIÓN */}
        <OrbitControls 
          makeDefault
          // 🎯 CAMBIO 2: Subimos el punto de mira a 0.8 metros. Esto hace que la cámara 
          // levante la cabeza, centrando el coche verticalmente en tu pantalla.
          target={[0, -10, 0]} 
          
          enablePan={false} 
          enableZoom={false} 
          
          // 🔒 CAMBIO 3: Damos un margen extra al suelo (1.65 radianes) para que al 
          // rotar de frente el límite invisible no te bloquee la cámara.
          maxPolarAngle={1.65} 
          minPolarAngle={0.2}
        />
      </Canvas>
    </div>
  );
}