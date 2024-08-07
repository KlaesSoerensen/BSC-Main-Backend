namespace ASP.NETCoreWebAPI.Models
{
    public class AssetCollection
    {
        public int Id { get; set; } // Primary Key
        public required string Name { get; set; } // Key Attribute - Required
        public bool Original { get; set; } // Attribute - Indicates if this is the original Asset Collection
        public CollectionEntry? Entries { get; set; } // Multi-Valued Attribute - Nullable, represents a collection of entries (multiple?)
        public AssetCollection? OriginalReference { get; set; } // Foreign Key (Self-Reference) - Nullable, references the original Asset Collection
    }

    // Need changes to self reference?
}
