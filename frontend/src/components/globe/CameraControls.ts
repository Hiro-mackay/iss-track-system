type CesiumModule = typeof import("cesium");
type Viewer = InstanceType<CesiumModule["Viewer"]>;
type Entity = InstanceType<CesiumModule["Entity"]>;

export function setCameraFollow(
  viewer: Viewer,
  entity: Entity,
  follow: boolean,
  onFollowChange: (follow: boolean) => void
): () => void {
  if (follow) {
    viewer.trackedEntity = entity;
  } else {
    viewer.trackedEntity = undefined;
  }

  const removeListener = viewer.trackedEntityChanged.addEventListener(
    (tracked: Entity | undefined) => {
      if (!tracked && follow) {
        onFollowChange(false);
      }
    }
  );

  return () => {
    removeListener();
  };
}

export function setInitialView(
  viewer: Viewer,
  Cesium: CesiumModule
): void {
  viewer.camera.setView({
    destination: Cesium.Cartesian3.fromDegrees(0, 20, 20_000_000),
  });
}
