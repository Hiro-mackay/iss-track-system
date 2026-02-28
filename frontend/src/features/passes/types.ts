export type Brightness = "bright" | "moderate" | "dim";

export interface PassEvent {
  time: string;
  azimuth: number;
  elevation: number;
}

export interface PassPrediction {
  rise: PassEvent;
  culmination: PassEvent;
  set: PassEvent;
  max_elevation: number;
  is_visible: boolean;
  brightness: Brightness;
}
