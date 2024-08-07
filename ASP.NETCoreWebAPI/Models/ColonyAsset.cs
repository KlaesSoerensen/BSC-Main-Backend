namespace ASP.NETCoreWebAPI.Models
{
    public class ColonyAsset
    {
        public int Id { get; set; } // Primary Key
        public AssetCollection AssetCollection { get; set; } // Foreign Key (AssetCollection)
        public Transform Transform { get; set; } // Foreign Key (Transform)
    }
}
