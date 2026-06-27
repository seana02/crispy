import { useParams } from "@solidjs/router";
import { createMemo, createResource, ErrorBoundary, Suspense } from "solid-js";
import TagForm from "src/components/TagForm";
import { getTagById, submitEditTag } from "src/stores/tagStore";


export default function EditTag() {
    const params = useParams();

    const id = +params.id!;

    const [tag] = createResource(() => id, () => getTagById(id));

    const data = createMemo(() => {
        const item = tag();
        if (!item) return null;
        return {
            id: item.id,
            name: item.name
        };
    });

    return (
        <ErrorBoundary fallback={err => <p>Loading error: {err.message}</p>}>
            <Suspense fallback={<p>Loading</p>}>
                <TagForm
                    id={id}
                    name={data()!.name}
                    submit={(name) => submitEditTag(id, name)}
                />
            </Suspense>
        </ErrorBoundary>
    );
}
