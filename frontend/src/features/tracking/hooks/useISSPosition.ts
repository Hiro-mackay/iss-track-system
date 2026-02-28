"use client";

import { useEffect, useRef, useState, useCallback } from "react";
import type { ISSPosition } from "../types";
import { getWebSocketUrl, fetchPosition, isISSPosition } from "../api";

const MAX_RECONNECT_DELAY = 30000;
const INITIAL_RECONNECT_DELAY = 1000;
const REST_POLL_INTERVAL = 2000;

export function useISSPosition() {
  const [position, setPosition] = useState<ISSPosition | null>(null);
  const [connected, setConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectDelayRef = useRef(INITIAL_RECONNECT_DELAY);
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pollTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const connectRef = useRef<() => void>(null);

  const stopPolling = useCallback(() => {
    if (pollTimerRef.current) {
      clearInterval(pollTimerRef.current);
      pollTimerRef.current = null;
    }
  }, []);

  const startPolling = useCallback(() => {
    stopPolling();
    const poll = async () => {
      try {
        const pos = await fetchPosition();
        setPosition(pos);
      } catch {
        // Silently fail, will retry
      }
    };
    poll();
    pollTimerRef.current = setInterval(poll, REST_POLL_INTERVAL);
  }, [stopPolling]);

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) return;

    const ws = new WebSocket(getWebSocketUrl());
    wsRef.current = ws;

    ws.onopen = () => {
      setConnected(true);
      reconnectDelayRef.current = INITIAL_RECONNECT_DELAY;
      stopPolling();
    };

    ws.onmessage = (event) => {
      try {
        const raw: unknown = JSON.parse(event.data);
        if (isISSPosition(raw)) {
          setPosition(raw);
        }
      } catch {
        // malformed frame, next tick will arrive
      }
    };

    ws.onclose = () => {
      setConnected(false);
      startPolling();
      const delay = reconnectDelayRef.current;
      reconnectDelayRef.current = Math.min(delay * 2, MAX_RECONNECT_DELAY);
      reconnectTimerRef.current = setTimeout(() => connectRef.current?.(), delay);
    };

    ws.onerror = () => {
      ws.close();
    };
  }, [startPolling, stopPolling]);

  useEffect(() => {
    connectRef.current = connect;
  }, [connect]);

  useEffect(() => {
    connect();
    return () => {
      if (reconnectTimerRef.current) clearTimeout(reconnectTimerRef.current);
      stopPolling();
      wsRef.current?.close();
    };
  }, [connect, stopPolling]);

  return { position, connected };
}
