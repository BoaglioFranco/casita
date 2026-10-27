import { useRef, useState } from 'react';
import type { MouseEvent, PointerEvent } from 'react';

type Gesture = {
  id: number;
  x: number;
  y: number;
  width: number;
  offset: number;
  phase: 'pending' | 'swiping' | 'ignored';
};

export function useSwipeToDelete({ disabled, onSwipeStart, onDelete }: {
  disabled: boolean;
  onSwipeStart: () => void;
  onDelete: () => void;
}) {
  const gesture = useRef<Gesture | null>(null);
  const suppressClick = useRef(false);
  const [offset, setOffset] = useState(0);
  const [dragging, setDragging] = useState(false);
  const [armed, setArmed] = useState(false);

  function finish(event: PointerEvent<HTMLElement>, cancelled: boolean) {
    const current = gesture.current;
    if (!current || current.id !== event.pointerId) return;
    // Clear before releasing capture, which can dispatch lostpointercapture.
    gesture.current = null;
    const shouldDelete = !cancelled && !disabled && current.phase === 'swiping'
      && current.offset >= current.width * 0.4;
    setOffset(0);
    setDragging(false);
    setArmed(false);
    if (event.currentTarget.hasPointerCapture(event.pointerId)) {
      event.currentTarget.releasePointerCapture(event.pointerId);
    }
    if (shouldDelete) onDelete();
  }

  return {
    offset,
    dragging,
    armed,
    handlers: {
      onPointerDown(event: PointerEvent<HTMLElement>) {
        if (event.pointerType !== 'touch' && event.pointerType !== 'mouse') return;
        if (event.button !== 0) return;
        suppressClick.current = false;
        if (disabled || !event.isPrimary || gesture.current
          || (event.target as Element).closest('[data-swipe-exclude]')) return;
        gesture.current = {
          id: event.pointerId, x: event.clientX, y: event.clientY,
          width: event.currentTarget.getBoundingClientRect().width,
          offset: 0, phase: 'pending',
        };
      },
      onPointerMove(event: PointerEvent<HTMLElement>) {
        const current = gesture.current;
        if (!current || current.id !== event.pointerId) return;
        if (disabled) { finish(event, true); return; }
        const dx = event.clientX - current.x;
        const dy = event.clientY - current.y;
        if (current.phase === 'pending') {
          if (Math.max(Math.abs(dx), Math.abs(dy)) < 10) return;
          suppressClick.current = true;
          if (dx >= 0 || -dx <= Math.abs(dy) * 1.5 || current.width <= 0) {
            current.phase = 'ignored';
            return;
          }
          current.phase = 'swiping';
          event.currentTarget.setPointerCapture(event.pointerId);
          setDragging(true);
          onSwipeStart();
        }
        if (current.phase !== 'swiping') return;
        current.offset = Math.min(current.width, Math.max(0, -dx));
        setOffset(current.offset);
        setArmed(current.offset >= current.width * 0.4);
      },
      onPointerUp(event: PointerEvent<HTMLElement>) { finish(event, false); },
      onPointerLeave(event: PointerEvent<HTMLElement>) {
        if (gesture.current?.phase !== 'swiping') finish(event, true);
      },
      onPointerCancel(event: PointerEvent<HTMLElement>) {
        if (gesture.current?.id === event.pointerId) suppressClick.current = true;
        finish(event, true);
      },
      onLostPointerCapture(event: PointerEvent<HTMLElement>) { finish(event, true); },
      onClickCapture(event: MouseEvent<HTMLElement>) {
        if (event.detail !== 0 && suppressClick.current) {
          event.preventDefault();
          event.stopPropagation();
          suppressClick.current = false;
        }
      },
    },
  };
}
