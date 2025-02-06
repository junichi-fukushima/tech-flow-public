import React, {useState} from 'react';
import {Button, Flex, Modal, Tag, message} from 'antd';
import {categories} from "@/util/menus";
import {postFavoritesCategories} from "@/api/favorites_categories";

const GetStarted: React.FC = () => {
  const [open, setOpen] = useState(true);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [selectedTags, setSelectedTags] = React.useState<string[]>([]);
  const [messageApi, contextHolder] = message.useMessage();

  const handleOk = async () => {
    setConfirmLoading(true);
    await postFavoritesCategories({liked_categories: selectedTags})
    messageApi.open({
      type: 'success',
      content: 'あなたにおすすめの記事を表示します',
    });
    setOpen(false);
    setConfirmLoading(false);
    // 画面をリロードする
    window.location.reload();
  };

  return (
    <>
      {contextHolder}
      <Modal
        title={
          <div style={{fontSize: 20}}>
            興味のある記事について教えてください
          </div>
        }
        open={open}
        onOk={handleOk}
        confirmLoading={confirmLoading}
        footer={[
          <Button
            key="submit"
            type="primary"
            loading={confirmLoading}
            onClick={handleOk}
            disabled={selectedTags.length === 0}
          >
            記事を見る
          </Button>,
        ]}
        closeIcon={null}
      >
        <div style={{minHeight: 180}}>
          <SelectTags
            selectedTags={selectedTags}
            setSelectedTags={setSelectedTags}
          />
        </div>
      </Modal>
    </>
  );
};

const SelectTags: React.FC<{
  selectedTags: string[], setSelectedTags: (tags: string[]) => void
}> = ({
        selectedTags,
        setSelectedTags
      }) => {
  const handleChange = (tag: string, checked: boolean) => {
    const nextSelectedTags = checked
      ? [...selectedTags, tag]
      : selectedTags.filter((t) => t !== tag);
    setSelectedTags(nextSelectedTags);
  };

  return (
    <Flex gap={4} wrap align="center">
      {categories.map<React.ReactNode>((category) => (
        <Tag.CheckableTag
          style={{fontSize: 12, padding: "8px"}}
          key={category}
          checked={selectedTags.includes(category)}
          onChange={(checked) => handleChange(category, checked)}
        >
          {category}
        </Tag.CheckableTag>
      ))}
    </Flex>
  );
};

export default GetStarted;
