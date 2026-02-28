"use client";

import { useRef, useEffect, useState } from "react";
import type { ISSPosition, OrbitPoint } from "@/lib/types";
import { configureCesiumBase, loadCesiumCSS, getViewerOptions } from "./cesium-setup";
import { createISSEntity, updateISSPosition } from "./ISSEntity";
import { replaceOrbitPolylines } from "./OrbitPath3D";
import { setCameraFollow, setInitialView } from "./CameraControls";

type CesiumModule = typeof import("cesium");
type Viewer = InstanceType<CesiumModule["Viewer"]>;
type Entity = InstanceType<CesiumModule["Entity"]>;

interface Props {
  position: ISSPosition | null;
  orbitPoints?: OrbitPoint[];
  followISS: boolean;
  onFollowChange: (follow: boolean) => void;
}

export default function GlobeView({
  position,
  orbitPoints = [],
  followISS,
  onFollowChange,
}: Props) {
  const containerRef = useRef<HTMLDivElement>(null);
  const viewerRef = useRef<Viewer | null>(null);
  const cesiumRef = useRef<CesiumModule | null>(null);
  const issRef = useRef<Entity | null>(null);
  const pastRef = useRef<Entity | null>(null);
  const futureRef = useRef<Entity | null>(null);
  const [viewerReady, setViewerReady] = useState(false);

  useEffect(() => {
    let destroyed = false;

    async function init() {
      if (!containerRef.current) return;
      configureCesiumBase();
      loadCesiumCSS();

      const Cesium = await import("cesium");
      if (destroyed) return;

      cesiumRef.current = Cesium;
      const viewer = new Cesium.Viewer(containerRef.current, getViewerOptions(Cesium));
      viewer.scene.globe.enableLighting = true;
      viewerRef.current = viewer;

      issRef.current = createISSEntity(Cesium, viewer);
      setInitialView(viewer, Cesium);
      setViewerReady(true);
    }

    void init().catch(() => {
      // init failed (e.g. Cesium load error); cleanup will run via destroyed flag
    });

    return () => {
      destroyed = true;
      if (viewerRef.current && !viewerRef.current.isDestroyed()) {
        viewerRef.current.destroy();
      }
      viewerRef.current = null;
      cesiumRef.current = null;
      issRef.current = null;
      pastRef.current = null;
      futureRef.current = null;
    };
  }, []);

  useEffect(() => {
    if (!viewerReady || !cesiumRef.current || !issRef.current || !position) return;
    updateISSPosition(cesiumRef.current, issRef.current, position);
  }, [viewerReady, position]);

  useEffect(() => {
    if (!viewerReady || !cesiumRef.current || !viewerRef.current) return;
    if (orbitPoints.length === 0) return;
    const { past, future } = replaceOrbitPolylines(
      cesiumRef.current,
      viewerRef.current,
      pastRef.current,
      futureRef.current,
      orbitPoints
    );
    pastRef.current = past;
    futureRef.current = future;
  }, [viewerReady, orbitPoints]);

  useEffect(() => {
    if (!viewerReady || !viewerRef.current || !issRef.current) return;
    return setCameraFollow(
      viewerRef.current,
      issRef.current,
      followISS,
      onFollowChange
    );
  }, [viewerReady, followISS, onFollowChange]);

  return <div ref={containerRef} className="h-full w-full" />;
}
