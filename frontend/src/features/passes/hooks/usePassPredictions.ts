"use client";

import useSWR from "swr";
import { fetchPasses } from "../api";
import type { PassPrediction } from "../types";

export function usePassPredictions(lat: number, lon: number, days = 3) {
  const { data, error, isLoading } = useSWR<PassPrediction[]>(
    lat && lon ? `passes-${lat}-${lon}-${days}` : null,
    () => fetchPasses(lat, lon, days),
    { refreshInterval: 300000 }
  );
  return { passes: data ?? [], error, isLoading };
}
