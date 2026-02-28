"use client";

import { useEffect, useState, useRef } from "react";

const DEBOUNCE_MS = 10_000;
const NOMINATIM_URL = "https://nominatim.openstreetmap.org/reverse";

export function useReverseGeocode(lat: number, lon: number) {
  const [location, setLocation] = useState("");
  const [loading, setLoading] = useState(false);
  const lastFetchRef = useRef(0);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    const controller = new AbortController();

    const fetchLocation = async () => {
      setLoading(true);
      try {
        const params = new URLSearchParams({
          format: "json",
          lat: lat.toString(),
          lon: lon.toString(),
          zoom: "5",
        });
        const res = await fetch(`${NOMINATIM_URL}?${params}`, {
          headers: { "User-Agent": "iss-tracker/1.0" },
          signal: controller.signal,
        });
        if (!res.ok) throw new Error("fetch failed");
        const data = await res.json();
        setLocation(data.display_name || "Ocean");
      } catch (err) {
        if (err instanceof Error && err.name === "AbortError") return;
        setLocation("Unknown location");
      } finally {
        setLoading(false);
        lastFetchRef.current = Date.now();
      }
    };

    if (timerRef.current) clearTimeout(timerRef.current);
    const elapsed = Date.now() - lastFetchRef.current;
    const delay = Math.max(0, DEBOUNCE_MS - elapsed);
    timerRef.current = setTimeout(fetchLocation, delay);

    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
      controller.abort();
    };
  }, [lat, lon]);

  return { location, loading };
}
