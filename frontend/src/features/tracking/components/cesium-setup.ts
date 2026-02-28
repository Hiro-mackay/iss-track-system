declare global {
  interface Window {
    CESIUM_BASE_URL: string;
  }
}

export function configureCesiumBase(): void {
  window.CESIUM_BASE_URL = "/cesium/";
}

export function loadCesiumCSS(): void {
  if (document.querySelector('link[href="/cesium/Widgets/widgets.css"]')) {
    return;
  }
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = "/cesium/Widgets/widgets.css";
  document.head.appendChild(link);
}

export function getViewerOptions(Cesium: typeof import("cesium")) {
  return {
    baseLayerPicker: false,
    geocoder: false,
    homeButton: false,
    infoBox: false,
    navigationHelpButton: false,
    sceneModePicker: false,
    selectionIndicator: false,
    timeline: false,
    animation: false,
    fullscreenButton: false,
    vrButton: false,
    baseLayer: new Cesium.ImageryLayer(
      new Cesium.UrlTemplateImageryProvider({
        url: "https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png",
        subdomains: "abcd",
        maximumLevel: 18,
        credit: new Cesium.Credit("CARTO"),
      })
    ),
  };
}
