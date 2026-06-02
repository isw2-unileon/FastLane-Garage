import { useGLTF } from '@react-three/drei';
import * as THREE from 'three';
import { useMemo, useEffect } from 'react';

interface CarModelProps {
  selectedZone: string | null;
  onSelectZone: (zone: string | null) => void;
}

export default function CarModel({ selectedZone, onSelectZone }: CarModelProps) {
  const { scene } = useGLTF('/src/assets/models/coche.glb');

  // 1. Configuración de los nuevos materiales con corrección de Z-Buffer (depthWrite)
  const materials = useMemo(() => {
    // Material base de superficie (cian semi-transparente)
    const surfaceMaterialCyan = new THREE.MeshStandardMaterial({
      color: new THREE.Color('#2379dc'),
      transparent: true,
      opacity: 0.4,
      blending: THREE.AdditiveBlending,
      wireframe: false,
      depthWrite: false,             // 🌟 Evita que las caras delanteras oculten las traseras
      side: THREE.DoubleSide,        // 🌟 Renderiza el interior y exterior de la chapa
    });

    // Material de wireframe base (cian) para superponer encima
    const wireframeMaterialCyan = new THREE.MeshBasicMaterial({
      color: new THREE.Color('#22d3ee'),
      wireframe: true,
      transparent: true,
      opacity: 0.10,                 // Más tenue para no saturar la vista
      depthWrite: false,             // 🌟 Evita parpadeos en las aristas traseras
    });

    // Material de superficie interactiva (naranja con más cuerpo/opacidad)
    const surfaceMaterialOrange = new THREE.MeshStandardMaterial({
      color: new THREE.Color('#ff5900'),
      transparent: true,
      opacity: 0.7,                  // Más sólido para que se distingan las ruedas
      blending: THREE.AdditiveBlending,
      wireframe: false,
      depthWrite: false,
      side: THREE.DoubleSide,
    });

    // Material de wireframe interactivo (naranja)
    const wireframeMaterialOrange = new THREE.MeshBasicMaterial({
      color: new THREE.Color('#eb280e'),
      wireframe: true,
      transparent: true,
      opacity: 0.8,
      depthWrite: false,
    });

    // Material de superficie seleccionada (naranja brillante de alerta)
    const selectedSurfaceMaterial = new THREE.MeshStandardMaterial({
      color: new THREE.Color('#f91616'),
      transparent: true,
      opacity: 0.8,
      wireframe: false,
      side: THREE.DoubleSide,        // Mantenemos DoubleSide para consistencia de click
    });

    return {
      surfaceMaterialCyan,
      wireframeMaterialCyan,
      surfaceMaterialOrange,
      wireframeMaterialOrange,
      selectedSurfaceMaterial,
    };
  }, []);

  // 2. Duplicación de mallas: Creamos grupos con Capa Sólida + Capa de Líneas
  const processedScene = useMemo(() => {
    const clone = scene.clone();
    const meshesToReplace: { child: THREE.Mesh; group: THREE.Group }[] = [];

    clone.traverse((child) => {
      if ((child as THREE.Mesh).isMesh) {
        const mesh = child as THREE.Mesh;
        const group = new THREE.Group();
        group.name = mesh.name;      // Preservamos el nombre para los eventos de click

        // Creamos los clones geométricos independientes
        const surfaceMesh = new THREE.Mesh(mesh.geometry);
        const wireframeMesh = new THREE.Mesh(mesh.geometry);

        // Asignamos los materiales según el sufijo interactivo del objeto
        if (mesh.name.toLowerCase().endsWith('_inter')) {
          surfaceMesh.material = materials.surfaceMaterialOrange;
          wireframeMesh.material = materials.wireframeMaterialOrange;
        } else {
          surfaceMesh.material = materials.surfaceMaterialCyan;
          wireframeMesh.material = materials.wireframeMaterialCyan;
        }

        group.add(surfaceMesh);
        group.add(wireframeMesh);
        meshesToReplace.push({ child: mesh, group: group });
      }
    });

    // Reemplazamos físicamente la estructura original en el árbol 3D
    meshesToReplace.forEach(({ child, group }) => {
      if (child.parent) {
        child.parent.add(group);
        child.parent.remove(child);
      }
    });

    return clone;
  }, [scene, materials]);

  // 3. Actualizador de estado visual dinámico al seleccionar una zona activa
  useEffect(() => {
    processedScene.traverse((child) => {
      if (child instanceof THREE.Group && child.name.toLowerCase().endsWith('_inter')) {
        const surfaceMesh = child.children[0] as THREE.Mesh;
        const wireframeMesh = child.children[1] as THREE.Mesh;

        if (selectedZone === child.name) {
          surfaceMesh.material = materials.selectedSurfaceMaterial;
          wireframeMesh.material = materials.wireframeMaterialOrange;
        } else {
          surfaceMesh.material = materials.surfaceMaterialOrange;
          wireframeMesh.material = materials.wireframeMaterialOrange;
        }
      }
    });
  }, [processedScene, selectedZone, materials]);

  return (
    <primitive
      object={processedScene}
      scale={1.6}
      position={[0, 0, 0]}
      onClick={(e: any) => {
        e.stopPropagation(); // Evita que el click atraviese el coche y toque el suelo
        
        // Comprobamos si lo que hemos tocado pertenece a un grupo interactivo
        if (e.object.parent instanceof THREE.Group && e.object.parent.name.toLowerCase().endsWith('_inter')) {
          onSelectZone(e.object.parent.name);
        } else {
          onSelectZone(null);
        }
      }}
    />
  );
}