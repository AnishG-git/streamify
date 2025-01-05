import React from 'react';

interface TypographySmallProps {
  children: React.ReactNode;
  className?: string;
}

export function TypographySmall({ children, className = '' }: TypographySmallProps) {
  return (
    <small className={`text-sm font-medium ${className}`}>{children}</small>
  );
}