import React from 'react';

interface TypographyH1Props {
    children: React.ReactNode;
    className?: string;
}

export function TypographyH1({ children, className = '' }: TypographyH1Props) {
    return (
        <h1 className={`scroll-m-20 text-4xl font-bold tracking-tight lg:text-5xl ${className}`}>
            {children}
        </h1>
    );
}