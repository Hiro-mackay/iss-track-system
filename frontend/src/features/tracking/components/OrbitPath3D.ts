import type { OrbitPoint } from "../types";

type CesiumModule = typeof import("cesium");
type Viewer = InstanceType<CesiumModule["Viewer"]>;
type Entity = InstanceType<CesiumModule["Entity"]>;

// Earth rotates ~0.2507 deg/min (360 deg / 1436.068 min sidereal day)
const EARTH_ROT_DEG_PER_MIN = 360.0 / 1436.068;

function projectToOrbitalRing(
  points: OrbitPoint[],
  midIndex: number
): OrbitPoint[] {
  return points.map((p, i) => {
    const dtMinutes = i - midIndex;
    return {
      lat: p.lat,
      lon: p.lon + dtMinutes * EARTH_ROT_DEG_PER_MIN,
      altitude_km: p.altitude_km,
    };
  });
}

function toPositions(
  Cesium: CesiumModule,
  points: OrbitPoint[]
): InstanceType<CesiumModule["Cartesian3"]>[] {
  return points.map((p) =>
    Cesium.Cartesian3.fromDegrees(p.lon, p.lat, p.altitude_km * 1000)
  );
}

function addPolyline(
  Cesium: CesiumModule,
  viewer: Viewer,
  positions: InstanceType<CesiumModule["Cartesian3"]>[],
  alpha: number,
  depthAlpha: number
): Entity {
  const gold = Cesium.Color.fromCssColorString("#FFD700");
  return viewer.entities.add({
    polyline: {
      positions,
      width: 2,
      arcType: Cesium.ArcType.NONE,
      material: new Cesium.ColorMaterialProperty(gold.withAlpha(alpha)),
      depthFailMaterial: new Cesium.ColorMaterialProperty(
        gold.withAlpha(depthAlpha)
      ),
    },
  });
}

export function replaceOrbitPolylines(
  Cesium: CesiumModule,
  viewer: Viewer,
  oldPast: Entity | null,
  oldFuture: Entity | null,
  points: OrbitPoint[]
): { past: Entity; future: Entity } {
  if (oldPast) viewer.entities.remove(oldPast);
  if (oldFuture) viewer.entities.remove(oldFuture);

  const mid = Math.floor(points.length / 2);
  const ring = projectToOrbitalRing(points, mid);
  // Close the ring: the API covers ~90min but orbital period is ~92.7min
  ring.push(ring[0]);

  const pastPositions = toPositions(Cesium, ring.slice(0, mid + 1));
  const futurePositions = toPositions(Cesium, ring.slice(mid));

  const past = addPolyline(Cesium, viewer, pastPositions, 0.8, 0.15);
  const future = addPolyline(Cesium, viewer, futurePositions, 0.3, 0.08);
  return { past, future };
}
