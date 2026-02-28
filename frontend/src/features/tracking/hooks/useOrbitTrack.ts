"use client";

import useSWR from "swr";
import { fetchOrbit } from "../api";
import type { OrbitPoint } from "../types";

export function useOrbitTrack(minutes = 90, refreshInterval = 60000) {
  const { data } = useSWR<OrbitPoint[]>(
    `orbit-${minutes}`,
    () => fetchOrbit(minutes),
    { refreshInterval }
  );
  return data ?? [];
}
