import React, {useEffect, useRef, useState} from "react";
import {Typography} from "antd";

const {Text} = Typography;

export const EllipsisText = ({text}: { text: string }) => {
  const textRef = useRef<HTMLDivElement>(null);
  const [textWidth, setTextWidth] = useState(0);

  useEffect(() => {
    const updateTextWidth = () => {
      if (textRef.current) {
        setTextWidth(textRef.current.offsetWidth);
      }
    };

    updateTextWidth(); // 初期表示時に実行
    window.addEventListener('resize', updateTextWidth); // ウィンドウサイズ変更時に実行

    return () => {
      window.removeEventListener('resize', updateTextWidth); // クリーンアップ
    };
  }, [text]);

  return (
    <div ref={textRef}>
      <Text type="secondary" ellipsis={true}
            style={{fontWeight: "initial", fontSize: 13, width: textWidth}}>{text}</Text>
    </div>
  );
};
