"use client";

import useSWR from "swr";
import { fetchStatus } from "../api";
import type { ISSStatus } from "../types";

const REFRESH_INTERVAL = 60000; // 60 seconds

export function useStationStatus() {
  const { data, error, isLoading } = useSWR<ISSStatus>(
    "iss-status",
    fetchStatus,
    { refreshInterval: REFRESH_INTERVAL }
  );
  return { status: data ?? null, error, isLoading };
}
