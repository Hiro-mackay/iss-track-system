"use client";

import useSWR from "swr";
import { fetchOrbit } from "@/lib/api";
import type { OrbitPoint } from "@/lib/types";

export function useOrbitTrack(minutes = 90, refreshInterval = 60000) {
  const { data } = useSWR<OrbitPoint[]>(
    `orbit-${minutes}`,
    () => fetchOrbit(minutes),
    { refreshInterval }
  );
  return data ?? [];
}
