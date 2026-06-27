import TagForm from "src/components/TagForm";
import { submitNewTag } from "src/stores/tagStore";

export default function NewTag() {
    return (
        <TagForm
            submit={submitNewTag}
        />
    );
}
