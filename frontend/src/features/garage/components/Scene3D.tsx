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
      
      <Canvas camera={{ position: [0, 10, 30], fov: 45 }}>
        {/* Iluminación general potente */}
        <ambientLight intensity={0.7} />
        <directionalLight position={[10, 15, 10]} intensity={1.5} castShadow />
        <directionalLight position={[-10, 10, -10]} intensity={0.6} />

        <Suspense fallback={null}>
          
          <Center>
            <CarModel selectedZone={selectedZone} onSelectZone={onSelectZone} />
            <MapaModel />
          </Center>
        </Suspense>

       
        <OrbitControls 
          makeDefault
          
          target={[0, -10, 0]} 
          
          enablePan={false} 
          enableZoom={false} 
          
          
          maxPolarAngle={1.65} 
          minPolarAngle={0.2}
        />
      </Canvas>
    </div>
  );
}