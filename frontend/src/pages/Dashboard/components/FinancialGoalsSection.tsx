import React, { useMemo } from 'react';
import { DebtItem } from '../types';
import './FinancialGoalsSection.css';

interface FinancialGoalsSectionProps {
  debtsByBank: DebtItem[];
}

export const FinancialGoalsSection: React.FC<FinancialGoalsSectionProps> = ({
  debtsByBank,
}) => {
  // Calculate total debt and percentages for each bank
  const { chartData, percentages } = useMemo(() => {
    // 1. Calculate total sum of all debts
    const totalDebt = debtsByBank.reduce((sum, debt) => sum + debt.amount, 0);
    const total = totalDebt;
    if (total === 0) {
      return {
        chartData: [],
        percentages: new Map<string, number>(),
      };
    }

    // 2. Calculate percentage for each bank
    const percentagesMap = new Map<string, number>();
    debtsByBank.forEach((debt) => {
      const percentage = (debt.amount / total) * 100;
      percentagesMap.set(debt.id, percentage);
    });

    // 3. Calculate donut chart segments based on percentages
    const outerRadius = 56;
    const innerRadius = 42; // For donut effect (56 - 14, where 14 is half of strokeWidth 28)
    const centerX = 91.5;
    const centerY = 91.5;
    const gapDegrees = 2; // Gap between segments in degrees
    
    // Calculate total gap space (one gap between each pair of segments)
    const totalGaps = debtsByBank.length * gapDegrees;
    // Available degrees for actual segments (excluding gaps)
    const availableDegrees = 360 - totalGaps;
    
    let currentAngle = -90; // Start from top (in degrees)

    const chartSegments = debtsByBank.map((debt) => {
      // Get percentage for this segment
      const percentage = percentagesMap.get(debt.id) || 0;
      
      // Calculate angle based on percentage of available degrees
      // This ensures visual proportion matches the percentage
      const segmentAngle = (percentage / 100) * availableDegrees;
      
      // Convert start angle to radians
      const startAngleRad = (currentAngle * Math.PI) / 180;
      // Convert end angle to radians (start + segment angle)
      const endAngleRad = ((currentAngle + segmentAngle) * Math.PI) / 180;
      
      // Calculate points for outer arc (donut outer edge)
      const outerStartX = centerX + outerRadius * Math.cos(startAngleRad);
      const outerStartY = centerY + outerRadius * Math.sin(startAngleRad);
      const outerEndX = centerX + outerRadius * Math.cos(endAngleRad);
      const outerEndY = centerY + outerRadius * Math.sin(endAngleRad);
      
      // Calculate points for inner arc (donut inner edge)
      const innerStartX = centerX + innerRadius * Math.cos(startAngleRad);
      const innerStartY = centerY + innerRadius * Math.sin(startAngleRad);
      const innerEndX = centerX + innerRadius * Math.cos(endAngleRad);
      const innerEndY = centerY + innerRadius * Math.sin(endAngleRad);
      
      // Determine if we need large arc flag (for arcs > 180 degrees)
      const largeArcFlag = segmentAngle > 180 ? 1 : 0;
      
      // Create SVG path for donut segment
      // Path: Move to outer start -> Arc along outer edge -> Line to inner end -> Arc along inner edge (reverse) -> Close
      const pathData = [
        `M ${outerStartX} ${outerStartY}`, // Move to outer start point
        `A ${outerRadius} ${outerRadius} 0 ${largeArcFlag} 1 ${outerEndX} ${outerEndY}`, // Outer arc (clockwise)
        `L ${innerEndX} ${innerEndY}`, // Line to inner end point
        `A ${innerRadius} ${innerRadius} 0 ${largeArcFlag} 0 ${innerStartX} ${innerStartY}`, // Inner arc (counter-clockwise)
        'Z', // Close path back to start
      ].join(' ');

      // Move to next segment position: current angle + segment angle + gap
      currentAngle += segmentAngle + gapDegrees;

      return {
        ...debt,
        pathData,
        percentage,
      };
    });

    return {
      chartData: chartSegments,
      percentages: percentagesMap,
    };
  }, [debtsByBank]);

  const formatAmount = (amount: number): string => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  return (
    <section className="financial-goals-section">
      <div className="financial-goals-section__content">
        <h2 className="financial-goals-section__title">Банковские задолженности</h2>
        
        <div className="financial-goals-section__chart-container">
          <div className="financial-goals-section__chart">
            <svg width="183" height="183" viewBox="0 0 183 183">
              {chartData.map((item) => (
                <path
                  key={item.id}
                  d={item.pathData}
                  fill={item.color.startsWith('#') ? item.color : undefined}
                  className={!item.color.startsWith('#') ? item.color : ''}
                />
              ))}
            </svg>
          </div>
          
          <ul className="financial-goals-section__legend">
            {debtsByBank.map((debt) => {
              const percentage = percentages.get(debt.id) || 0;
              return (
                <li key={debt.id} className="financial-goals-section__legend-item">
                  <div
                    className="financial-goals-section__legend-dot"
                    style={
                      debt.color.startsWith('#')
                        ? { backgroundColor: debt.color }
                        : undefined
                    }
                  />
                  <div className="financial-goals-section__legend-info">
                    <span className="financial-goals-section__legend-bank">{debt.bankName}</span>
                    <span className="financial-goals-section__legend-percentage">
                      {percentage.toFixed(1)}%
                    </span>
                    <span className="financial-goals-section__legend-amount">
                      {formatAmount(debt.amount)}
                    </span>
                  </div>
                </li>
              );
            })}
          </ul>
        </div>
      </div>
    </section>
  );
};

