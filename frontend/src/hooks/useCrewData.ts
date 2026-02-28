"use client";

import useSWR from "swr";
import { fetchCrew } from "@/lib/api";
import { FALLBACK_CREW } from "@/lib/fallback-crew";
import type { CrewMember } from "@/lib/types";

const REFRESH_INTERVAL = 3600000; // 1 hour

export function useCrewData() {
  const { data, error, isLoading } = useSWR<CrewMember[]>(
    "iss-crew",
    fetchCrew,
    { refreshInterval: REFRESH_INTERVAL }
  );
  return {
    crew: data ?? (error ? FALLBACK_CREW : []),
    error,
    isLoading,
    isOffline: !!error,
  };
}
