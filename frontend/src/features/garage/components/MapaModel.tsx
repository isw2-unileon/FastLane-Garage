import { useGLTF } from '@react-three/drei';
import * as THREE from 'three';

export default function MapaModel() {
  // Cargamos tu mapa usando la ruta estática directa
  const { scene } = useGLTF('/src/assets/models/escenario.glb');

  // Recorremos el escenario para activar sus texturas, colores y emisiones
  scene.traverse((child) => {
    if ((child as THREE.Mesh).isMesh) {
      const mesh = child as THREE.Mesh;
      
      if (mesh.material instanceof THREE.MeshStandardMaterial) {
        mesh.material.needsUpdate = true;
        
        // Si tiene el mapa de emisión configurado en Blender, mantenemos el brillo activo
        if (mesh.material.emissive && mesh.material.emissive.getHex() !== 0x000000) {
          mesh.material.emissiveIntensity = 2.0; 
        }
      }
    }
  });

  return (
    <primitive 
      object={scene} 
      scale={1.0}        
      position={[0, -0.6, 0]} // Ajustado para que encaje con el suelo físico de la cámara
    />
  );
}