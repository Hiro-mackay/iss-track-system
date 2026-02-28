import type { ISSPosition } from "../types";

type CesiumModule = typeof import("cesium");

export function createISSEntity(
  Cesium: CesiumModule,
  viewer: InstanceType<CesiumModule["Viewer"]>
): InstanceType<CesiumModule["Entity"]> {
  const gold = Cesium.Color.fromCssColorString("#FFD700");

  return viewer.entities.add({
    name: "ISS",
    point: {
      pixelSize: 12,
      color: gold,
      outlineColor: Cesium.Color.WHITE,
      outlineWidth: 2,
    },
    label: {
      text: "ISS",
      font: "12px sans-serif",
      fillColor: gold,
      style: Cesium.LabelStyle.FILL_AND_OUTLINE,
      outlineWidth: 1,
      outlineColor: Cesium.Color.BLACK,
      verticalOrigin: Cesium.VerticalOrigin.BOTTOM,
      pixelOffset: new Cesium.Cartesian2(0, -16),
    },
  });
}

export function updateISSPosition(
  Cesium: CesiumModule,
  entity: InstanceType<CesiumModule["Entity"]>,
  position: ISSPosition
): void {
  const cartesian = Cesium.Cartesian3.fromDegrees(
    position.lon,
    position.lat,
    position.altitude_km * 1000
  );
  // ConstantPositionProperty implements PositionProperty but the dynamically-
  // imported Cesium module types do not satisfy the structural check.
  // This cast is safe: we are assigning a concrete PositionProperty implementor.
  entity.position = new Cesium.ConstantPositionProperty(
    cartesian
  ) as unknown as InstanceType<CesiumModule["Entity"]>["position"];
}
